#!/usr/bin/env bash
set -ex

mode="mainnet"

CONFIG_DIR=${HOME}/.fnsap
if [[ -z "${CHAIN_DIR}" ]]
then
  CHAIN_DIR=${CONFIG_DIR}
fi

if [[ $1 == "docker" ]]
then
    if [[ $2 == "testnet" ]]
    then
        mode="testnet"
    fi
    FNSAD="docker run -i -p 26656:26656 -p 26657:26657 -v ${CONFIG_DIR}:/root/.fnsap finschia/finschianode-proxy fnsad-proxy"
    CHAIN_DIR="/root/.fnsap"
elif [[ $1 == "testnet" ]]
then
    mode="testnet"
fi

FNSAD=${FNSAD:-fnsad-proxy}

# initialize
rm -rf ${CONFIG_DIR}

# TODO
# Configure your CLI to eliminate need for chain-id flag
#${FNSAD} config chain-id finschia
#${FNSAD} config output json
#${FNSAD} config indent true
#${FNSAD} config trust-node true
#${FNSAD} config keyring-backend test

# Initialize configuration files and genesis file
# moniker is the name of your node
${FNSAD} init solo --chain-id=finschia --home=${CHAIN_DIR}

# configure for testnet
if [[ ${mode} == "testnet" ]]
then
    if [[ $1 == "docker" ]]
    then
        docker run -i -p 26656:26656 -p 26657:26657 -v ${CONFIG_DIR}:/root/.fnsap finschia/finschianode-proxy sh -c "export FNSAD_TESTNET=true"
    else
       export FNSAD_TESTNET=true
    fi
fi

# Please do not use the TEST_MNEMONIC for production purpose
TEST_MNEMONIC="mind flame tobacco sense move hammer drift crime ring globe art gaze cinnamon helmet cruise special produce notable negative wait path scrap recall have"

${FNSAD} keys add jack --home=${CHAIN_DIR} --keyring-backend=test --recover --account=0 <<< ${TEST_MNEMONIC}
${FNSAD} keys add alice --home=${CHAIN_DIR} --keyring-backend=test --recover --account=1 <<< ${TEST_MNEMONIC}
${FNSAD} keys add bob --home=${CHAIN_DIR} --keyring-backend=test --recover --account=2 <<< ${TEST_MNEMONIC}
${FNSAD} keys add rinah --home=${CHAIN_DIR} --keyring-backend=test --recover --account=3 <<< ${TEST_MNEMONIC}
${FNSAD} keys add sam --home=${CHAIN_DIR} --keyring-backend=test --recover --account=4 <<< ${TEST_MNEMONIC}
${FNSAD} keys add evelyn --home=${CHAIN_DIR} --keyring-backend=test --recover --account=5 <<< ${TEST_MNEMONIC}

# TODO
#if [[ ${mode} == "testnet" ]]
#then
#   ${FNSAD} add-genesis-account tlink15la35q37j2dcg427kfy4el2l0r227xwhc2v3lg 9223372036854775807link,1stake
#else
#   ${FNSAD} add-genesis-account link15la35q37j2dcg427kfy4el2l0r227xwhuaapxd 9223372036854775807link,1stake
#fi
# Add both accounts, with coins to the genesis file
${FNSAD} add-genesis-account $(${FNSAD} keys show jack -a --home=${CHAIN_DIR} --keyring-backend=test) 1000link,1000000000000stake --home=${CHAIN_DIR}
${FNSAD} add-genesis-account $(${FNSAD} keys show alice -a --home=${CHAIN_DIR} --keyring-backend=test) 1000link,1000000000000stake --home=${CHAIN_DIR}
${FNSAD} add-genesis-account $(${FNSAD} keys show bob -a --home=${CHAIN_DIR} --keyring-backend=test) 1000link,1000000000000stake --home=${CHAIN_DIR}
${FNSAD} add-genesis-account $(${FNSAD} keys show rinah -a --home=${CHAIN_DIR} --keyring-backend=test) 1000link,1000000000000stake --home=${CHAIN_DIR}
${FNSAD} add-genesis-account $(${FNSAD} keys show sam -a --home=${CHAIN_DIR} --keyring-backend=test) 1000link,1000000000000stake --home=${CHAIN_DIR}
${FNSAD} add-genesis-account $(${FNSAD} keys show evelyn -a --home=${CHAIN_DIR} --keyring-backend=test) 1000link,1000000000000stake --home=${CHAIN_DIR}

${FNSAD} gentx jack 100000000stake --home=${CHAIN_DIR} --keyring-backend=test --chain-id=finschia

${FNSAD} collect-gentxs --home=${CHAIN_DIR}

${FNSAD} validate-genesis --home=${CHAIN_DIR}

# ${FNSAD} start --log_level *:debug --rpc.laddr=tcp://0.0.0.0:26657 --p2p.laddr=tcp://0.0.0.0:26656