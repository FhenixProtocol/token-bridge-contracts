/*
 * Copyright 2019, Offchain Labs, Inc.
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

#ifndef checkpointutils_hpp
#define checkpointutils_hpp

#include <avm_values/codepoint.hpp>
#include <avm_values/tuple.hpp>

struct MachineStateKeys {
    uint256_t register_hash;
    uint256_t datastack_hash;
    uint256_t auxstack_hash;
    CodePointStub pc;
    CodePointStub err_pc;
    unsigned char status_char;
};

extern std::unordered_map<int, int> blockreason_type_length;

namespace checkpoint {
namespace utils {
std::vector<unsigned char> serializeValue(const value& val);
CodePointStub deserializeCodePointStub(const char*& bufptr);
uint256_t deserializeUint256_t(const std::vector<unsigned char>& val);
std::vector<std::vector<unsigned char>> parseTuple(
    const std::vector<unsigned char>& data);
MachineStateKeys extractStateKeys(
    const std::vector<unsigned char>& stored_state);
std::vector<unsigned char> serializeStateKeys(
    const MachineStateKeys& state_data);
std::vector<unsigned char> GetHashKey(const value& val);
}  // namespace utils
}  // namespace checkpoint

#endif /* checkpointutils_hpp */
