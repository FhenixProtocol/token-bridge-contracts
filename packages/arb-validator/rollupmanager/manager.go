/*
* Copyright 2020, Offchain Labs, Inc.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package rollupmanager

import (
	"context"
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/offchainlabs/arbitrum/packages/arb-validator/structures"

	"github.com/offchainlabs/arbitrum/packages/arb-util/protocol"
	"github.com/offchainlabs/arbitrum/packages/arb-util/value"

	"github.com/offchainlabs/arbitrum/packages/arb-validator/rollup"

	"github.com/offchainlabs/arbitrum/packages/arb-util/common"

	"github.com/offchainlabs/arbitrum/packages/arb-validator/arbbridge"
)

const (
	maxReorgDepth = 100
)

type Manager struct {
	sync.Mutex
	RollupAddress common.Address
	client        arbbridge.ArbClient
	listeners     []rollup.ChainListener

	listenerAddChan chan rollup.ChainListener
	actionChan      chan func(*rollup.ChainObserver)
}

func CreateManager(
	ctx context.Context,
	rollupAddr common.Address,
	arbitrumCodeFilePath string,
	updateOpinion bool,
	clnt arbbridge.ArbClient,
	forceFreshStart bool,
	dbPrefix string,
	stressTest bool, // if true, generate artificial chaos to stress-test the implementation
) (*Manager, error) {
	man := &Manager{
		RollupAddress:   rollupAddr,
		client:          clnt,
		listenerAddChan: make(chan rollup.ChainListener, 10),
		actionChan:      make(chan func(*rollup.ChainObserver), 10),
	}
	go func() {
		for {
			runCtx, cancelFunc := context.WithCancel(ctx)

			checkpointer := rollup.NewProductionCheckpointer(
				runCtx,
				rollupAddr,
				arbitrumCodeFilePath,
				big.NewInt(maxReorgDepth),
				dbPrefix,
				forceFreshStart,
			)

			latestBlockId, chainObserverBuf, restoreCtx, err := checkpointer.RestoreLatestState(runCtx, clnt, rollupAddr, updateOpinion)
			if err != nil {
				log.Fatal(err)
			}
			watcher, err := clnt.NewRollupWatcher(rollupAddr)
			if err != nil {
				log.Fatal(err)
			}
			if stressTest {
				watcher = NewStressTestWatcher(watcher, 30*time.Second)
			}
			chain := chainObserverBuf.UnmarshalFromCheckpoint(runCtx, restoreCtx, latestBlockId, watcher, checkpointer)

			man.Lock()
			// Clear pending listeners
			for len(man.listenerAddChan) > 0 {
				<-man.listenerAddChan
			}
			// Add manager's listeners
			for _, listener := range man.listeners {
				chain.AddListener(listener)
			}
			man.Unlock()

			chain.Start(runCtx)

			current, err := clnt.CurrentBlockId(runCtx)
			if err != nil {
				log.Fatal(err)
			}

			headersChan, err := clnt.SubscribeBlockHeaders(runCtx, latestBlockId)
			if err != nil {
				blockId, err := clnt.BlockIdForHeight(ctx, common.NewTimeBlocks(big.NewInt(0)))
				if err != nil {
					panic(err)
				}
				log.Println("Error subscribing to block headers", latestBlockId.HeaderHash, latestBlockId.Height.AsInt(), blockId.HeaderHash, blockId.Height.AsInt(), err)

				cancelFunc()
				time.Sleep(2 * time.Second)
				continue
			}
			reachedHead := false
		runLoop:
			for {
				select {
				case maybeBlockId, ok := <-headersChan:
					if !ok {
						log.Println("Manager stopped receiving headers")
						break runLoop
					}
					if maybeBlockId.Err != nil {
						log.Println("Error getting new header", maybeBlockId.Err)
						break runLoop
					}

					blockId := maybeBlockId.BlockId

					if !reachedHead && blockId.Height.Cmp(current.Height) >= 0 {
						log.Println("Reached head")
						reachedHead = true
						chain.NowAtHead()
						log.Println("Now at head")
					}

					chain.NotifyNewBlock(blockId.Clone())

					events, err := watcher.GetEvents(runCtx, blockId)
					if err != nil {
						log.Println("Manager hit error getting events", err)
						break runLoop
					}
					for _, event := range events {
						handleNotification(runCtx, event, chain)
					}
				case action := <-man.actionChan:
					action(chain)
				}
			}

			cancelFunc()

			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(5 * time.Second) // give time for things to settle, post-reorg, before restarting stuff
			}
		}
	}()

	return man, nil
}

func (man *Manager) AddListener(listener rollup.ChainListener) {
	man.Lock()
	man.listeners = append(man.listeners, listener)
	man.listenerAddChan <- listener
	man.Unlock()
}

func (man *Manager) ExecuteCall(messages value.TupleValue, maxSteps uint32) (*protocol.ExecutionAssertion, uint32) {
	retChan := make(chan struct {
		*protocol.ExecutionAssertion
		uint32
	}, 1)
	man.actionChan <- func(chain *rollup.ChainObserver) {
		mach := chain.LatestKnownValidMachine()
		latestTime := chain.CurrentBlockId().Height
		timeBounds := &protocol.TimeBoundsBlocks{latestTime, latestTime}
		go func() {
			assertion, numSteps := mach.ExecuteAssertion(maxSteps, timeBounds, messages)
			retChan <- struct {
				*protocol.ExecutionAssertion
				uint32
			}{assertion, numSteps}
		}()
	}
	ret := <-retChan
	return ret.ExecutionAssertion, ret.uint32
}

func (man *Manager) CurrentBlockId() *structures.BlockId {
	retChan := make(chan *structures.BlockId, 1)
	man.actionChan <- func(chain *rollup.ChainObserver) {
		retChan <- chain.CurrentBlockId()
	}
	return <-retChan
}

func handleNotification(ctx context.Context, event arbbridge.Event, chain *rollup.ChainObserver) {
	chain.Lock()
	defer chain.Unlock()
	switch ev := event.(type) {
	case arbbridge.MessageDeliveredEvent:
		chain.MessageDelivered(ctx, ev)
	case arbbridge.StakeCreatedEvent:
		chain.CreateStake(ctx, ev)
	case arbbridge.ChallengeStartedEvent:
		chain.NewChallenge(ctx, ev)
	case arbbridge.ChallengeCompletedEvent:
		chain.ChallengeResolved(ctx, ev)
	case arbbridge.StakeRefundedEvent:
		chain.RemoveStake(ctx, ev)
	case arbbridge.PrunedEvent:
		chain.PruneLeaf(ctx, ev)
	case arbbridge.StakeMovedEvent:
		chain.MoveStake(ctx, ev)
	case arbbridge.AssertedEvent:
		err := chain.NotifyAssert(ctx, ev)
		if err != nil {
			panic(err)
		}
	case arbbridge.ConfirmedEvent:
		chain.ConfirmNode(ctx, ev)
	}
}
