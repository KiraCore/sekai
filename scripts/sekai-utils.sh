#!/usr/bin/env bash

function sekaiUtilsVersion() {
    sekaiUtilsSetup "version" 2> /dev/null || sekai-utils sekaiUtilsSetup "version"
}

# this is default installation script for sekaid utils
# ./sekai-utils.sh sekaiUtilsSetup
function sekaiUtilsSetup() {
    local SEKAI_UTILS_VERSION="v0.0.1.3"
    if [ "$1" == "version" ] ; then
        echo "$SEKAI_UTILS_VERSION"
        return 0
    else
        local UTILS_SOURCE=$(realpath "$0")
        local VERSION=$($UTILS_SOURCE sekaiUtilsVersion || echo '')
        local UTILS_DESTINATION="/usr/local/bin/sekai-utils.sh"

        if [ "$VERSION" != "$SEKAI_UTILS_VERSION" ] ; then
            bash-utils echoErr "ERROR: Self check version mismatch, expected '$SEKAI_UTILS_VERSION', but got '$VERSION'"
            return 1
        elif [ "$UTILS_SOURCE" == "$UTILS_DESTINATION" ] ; then
            bash-utils echoErr "ERROR: Installation source script and destination can't be the same"
            return 1
        elif [ ! -f $UTILS_SOURCE ] ; then
            bash-utils echoErr "ERROR: utils source was NOT found"
            return 1
        else
            bash-utils echoInfo "INFO: Utils source found"
            mkdir -p "/usr/local/bin"
            cp -fv "$UTILS_SOURCE" "$UTILS_DESTINATION"
            cp -fv "$UTILS_SOURCE" "/usr/local/bin/sekai-utils"
            chmod -v 555 $UTILS_DESTINATION "/usr/local/bin/sekai-utils"

            local SUDOUSER="${SUDO_USER}" && [ "$SUDOUSER" == "root" ] && SUDOUSER=""
            local USERNAME="${USER}" && [ "$USERNAME" == "root" ] && USERNAME=""
            local LOGNAME=$(logname 2> /dev/null echo "") && [ "$LOGNAME" == "root" ] && LOGNAME=""

            local TARGET="/$LOGNAME/.bashrc" && [ -f $TARGET ] && chmod 777 $TARGET && bash-utils echoInfo "INFO: /etc/profile executable target set to $TARGET"
            TARGET="/$USERNAME/.bashrc" && [ -f $TARGET ] && chmod 777 $TARGET && bash-utils echoInfo "INFO: /etc/profile executable target set to $TARGET"
            TARGET="/$SUDOUSER/.bashrc" && [ -f $TARGET ] && chmod 777 $TARGET && bash-utils echoInfo "INFO: /etc/profile executable target set to $TARGET"
            TARGET="/root/.bashrc" && [ -f $TARGET ] && chmod 777 $TARGET && bash-utils echoInfo "INFO: /etc/profile executable target set to $TARGET"
            TARGET=~/.bashrc && [ -f $TARGET ] && chmod 777 $TARGET && bash-utils echoInfo "INFO: /etc/profile executable target set to $TARGET"
            TARGET=~/.zshrc && [ -f $TARGET ] && chmod 777 $TARGET && bash-utils echoInfo "INFO: /etc/profile executable target set to $TARGET"
            TARGET=~/.profile && [ -f $TARGET ] && chmod 777 $TARGET && bash-utils echoInfo "INFO: /etc/profile executable target set to $TARGET"

            bash-utils setGlobEnv SEKAI_TOOLS_SRC "$UTILS_DESTINATION"

            local AUTOLOAD_SET=$(bash-utils getLastLineByPrefix "source $UTILS_DESTINATION" /etc/profile 2> /dev/null || echo "-1")

            if [[ $AUTOLOAD_SET -lt 0 ]] ; then
                echo "source $UTILS_DESTINATION || echo \"ERROR: Failed to load sekaid utils from '$UTILS_DESTINATION'\"" >> /etc/profile
            fi

            bash-utils loadGlobEnvs
            bash-utils echoInfo "INFO: SUCCESS, Installed sekai-utils $(sekai-utils sekaiUtilsVersion)"
        fi
    fi
}

function txQuery() {
    (! $(isTxHash "$1")) && echoErr "ERROR: Infalid Transaction Hash '$1'" && sekaid query tx "$1" --output=json --home=$SEKAID_HOME | jq || echoErr "ERROR: Transaction '$1' was NOT found or failed"
}

function txAwait() {
    local START_TIME="$(date -u +%s)"
    local RAW=""
    local TIMEOUT=""

    if (! $(isTxHash "$1")) ; then
        RAW=$(cat)
        TIMEOUT=$1
    else
        RAW=$1
        TIMEOUT=$2
    fi

    (! $(isNaturalNumber $TIMEOUT)) && TIMEOUT=0
    [[ $TIMEOUT -le 0 ]] && MAX_TIME="∞" || MAX_TIME="$TIMEOUT"
    
    local TXHASH=""
    if (! $(isTxHash "$RAW")) ; then
        # INPUT example: {"height":"0","txhash":"DF8BFCC9730FDBD33AEA184EC3D6C37B4311BC1C0E2296893BC020E4638A0D6F","codespace":"","code":0,"data":"","raw_log":"","logs":[],"info":"","gas_wanted":"0","gas_used":"0","tx":null,"timestamp":""}
        local VAL=$(echo $RAW | jsonParse "" 2> /dev/null || echo "")
        if ($(isNullOrEmpty "$VAL")) ; then
            echoErr "ERROR: Failed to propagate transaction:"
            echoErr "$RAW"
            return 1
        fi

        TXHASH=$(echo $VAL | jsonQuickParse "txhash" 2> /dev/null || echo "")
        if ($(isNullOrEmpty "$TXHASH")) ; then
            echoErr "ERROR: Transaction hash 'txhash' was NOT found in the tx propagation response:"
            echoErr "$RAW"
            return 1
        fi
    else
        TXHASH="${RAW^^}"
    fi

    echoInfo "INFO: Transaction hash '$TXHASH' was found!"
    echoInfo "INFO: Please wait for tx confirmation, timeout will occur in $MAX_TIME seconds ..."

    while : ; do
        local ELAPSED=$(($(date -u +%s) - $START_TIME))
        local OUT=$(sekaid query tx $TXHASH --output=json --home=$SEKAID_HOME 2> /dev/null | jsonParse "" 2> /dev/null || echo -n "")
        if [ ! -z "$OUT" ] ; then
            echoInfo "INFO: Transaction query response received received:"
            echo $OUT | jq

            local CODE=$(echo $OUT | jsonQuickParse "code" 2> /dev/null || echo -n "")
            if [ "$CODE" == "0" ] ; then
                echoInfo "INFO: Transaction was confirmed sucessfully!"
                return 0
            else
                echoErr "ERROR: Transaction failed with exit code '$CODE'"
                return 1
            fi
        else
            echoWarn "WAITING: Transaction is NOT confirmed yet, elapsed ${ELAPSED}/${MAX_TIME} s"
        fi

        if [[ $TIMEOUT -gt 0 ]] && [[ $ELAPSED -gt $TIMEOUT ]] ; then
            echoInfo "INFO: Transaction query response was NOT received:"
            echo $RAW | jq 2> /dev/null || echoErr "$RAW"
            echoErr "ERROR: Timeout, failed to confirm tx hash '$TXHASH' within ${TIMEOUT} s limit"
            return 1
        else
            sleep 0.5
        fi
    done
}

function txAwait2() {
    local START_TIME="$(date -u +%s)"
    local RAW=""
    local TIMEOUT=""

    if (! $(isTxHash "$1")) ; then
        RAW=$(cat)
        TIMEOUT=$1
    else
        RAW=$1
        TIMEOUT=$2
    fi

    (! $(isNaturalNumber $TIMEOUT)) && TIMEOUT=0
    [[ $TIMEOUT -le 0 ]] && MAX_TIME="∞" || MAX_TIME="$TIMEOUT"

    local TXHASH=""
    if (! $(isTxHash "$RAW")) ; then
        # INPUT example: {"height":"0","txhash":"DF8BFCC9730FDBD33AEA184EC3D6C37B4311BC1C0E2296893BC020E4638A0D6F","codespace":"","code":0,"data":"","raw_log":"","logs":[],"info":"","gas_wanted":"0","gas_used":"0","tx":null,"timestamp":""}
        local VAL=$(echo $RAW | jsonParse "" 2> /dev/null || echo "")
        if ($(isNullOrEmpty "$VAL")) ; then
            echoErr "ERROR: Failed to propagate transaction:"
            echoErr "$RAW"
            return 1
        fi

        TXHASH=$(echo $VAL | jsonQuickParse "txhash" 2> /dev/null || echo "")
        if ($(isNullOrEmpty "$TXHASH")) ; then
            echoErr "ERROR: Transaction hash 'txhash' was NOT found in the tx propagation response:"
            echoErr "$RAW"
            return 1
        fi
    else
        TXHASH="${RAW^^}"
    fi

    while : ; do
        local ELAPSED=$(($(date -u +%s) - $START_TIME))
        local OUT=$(sekaid query tx $TXHASH --output=json --home=$SEKAID_HOME 2> /dev/null | jsonParse "" 2> /dev/null || echo -n "")
        if [ ! -z "$OUT" ] ; then

            local CODE=$(echo $OUT | jsonQuickParse "code" 2> /dev/null || echo -n "")
            if [ "$CODE" == "0" ] ; then
                echo "${TXHASH,,}"
                return 0
            else
                echoErr "ERROR: Transaction failed with exit code '$CODE'"
                return 1
            fi
        fi

        if [[ $TIMEOUT -gt 0 ]] && [[ $ELAPSED -gt $TIMEOUT ]] ; then
            echoInfo "INFO: Transaction query response was NOT received:"
            echo $RAW | jq 2> /dev/null || echoErr "$RAW"
            echoErr "ERROR: Timeout, failed to confirm tx hash '$TXHASH' within ${TIMEOUT} s limit"
            return 1
        else
            sleep 0.5
        fi
    done

    echo "${TXHASH,,}"
}

# e.g. showAddress validator
function showAddress() {
    ($(isKiraAddress "$1")) && echo "$1" || echo $(sekaid keys show "$1" --keyring-backend=test --output=json --home=$SEKAID_HOME 2> /dev/null | jsonParse "address" 2> /dev/null || echo -n "")
}

function showKeys() {
    sekaid keys list --keyring-backend=test --output=json --home=$SEKAID_HOME | jq
}

# showPermissions validator
function showPermissions() {
    local ADDRESS=$(showAddress $1)
    echo $(sekaid query customgov permissions "$ADDRESS" --output=json --home=$SEKAID_HOME 2> /dev/null | jsonParse 2> /dev/null || echo -n "") && echo -n ""
}

function isPermBlacklisted() {
    local ADDR=$(showAddress "$1")
    local PERM=$2
    if (! $(isNaturalNumber $PERM)) || ($(isNullOrEmpty $ADDR)) ; then
        echo "false"
    else
        INDEX=$(showPermissions $ADDR 2> /dev/null | jq ".blacklist | index($PERM)" 2> /dev/null || echo -n "")
        ($(isNaturalNumber $INDEX)) && echo "true" || echo "false"
    fi
}

function isPermWhitelisted() {
    local ADDR=$(showAddress "$1")
    local PERM=$2
    if (! $(isNaturalNumber $PERM)) || ($(isNullOrEmpty $ADDR)) ; then
        echo "false"
    else
        INDEX=$(showPermissions $ADDR 2> /dev/null | jq ".whitelist | index($PERM)" 2> /dev/null || echo -n "")
        if ($(isNaturalNumber $INDEX)) && (! $(isPermBlacklisted $ADDR $PERM)) ; then
            echo "true" 
        else
            echo "false"
        fi
    fi
}

function lastProposal() {
    local BOTTOM_PROPOSALS=$(sekaid query customgov proposals --limit=1 --reverse --output=json --home=$SEKAID_HOME | jq -cr '.proposals | last | .proposal_id' 2> /dev/null || echo "")
    local TOP_PROPOSALS=$(sekaid query customgov proposals --limit=1 --output=json --home=$SEKAID_HOME | jq -cr '.proposals | last | .proposal_id' 2> /dev/null || echo "")
    ($(isNullOrEmpty "$BOTTOM_PROPOSALS")) && ($(isNullOrEmpty "$TOP_PROPOSALS")) && echo 0 && return 1
    (! $(isNaturalNumber $BOTTOM_PROPOSALS)) && (! $(isNaturalNumber $TOP_PROPOSALS)) && echo 0 && return 2
    [[ $TOP_PROPOSALS -le 0 ]] && [[ $BOTTOM_PROPOSALS -le 0 ]] && echo 0 && return 3
    [[ $TOP_PROPOSALS -gt $BOTTOM_PROPOSALS ]] && echo "$TOP_PROPOSALS" || echo "$BOTTOM_PROPOSALS"
    return 0
}

# voteProposal validator $(lastProposal) 0
function voteProposal() {
    local ACCOUNT=$1
    local PROPOSAL=$2
    local VOTE=$3
    
    echoInfo "INFO: Voting '$VOTE' on proposal '$PROPOSAL' with account '$ACCOUNT'"
    sekaid tx customgov proposal vote $PROPOSAL $VOTE --from=$ACCOUNT --chain-id=$NETWORK_NAME --keyring-backend=test  --fees=100ukex --yes --log_format=json --broadcast-mode=async --output=json --home=$SEKAID_HOME | txAwait
}

# voteYes $(lastProposal) validator
function voteYes() {
    voteProposal "$2" "$1" "1"
}

# voteNo $(lastProposal) validator
function voteNo() {
    voteProposal "$2" "$1" "0"
}

function showNetworkProperties() {
    local NETWORK_PROPERTIES=$(sekaid query customgov network-properties --output=json --home=$SEKAID_HOME 2> /dev/null || echo "" | jq -rc 2> /dev/null || echo "")
    ($(isNullOrEmpty "$NETWORK_PROPERTIES")) && echo -n "" && return 1
    echo $NETWORK_PROPERTIES
    return 0
}

# showVotes $(lastProposal) 
function showVotes() {
    sekaid query customgov votes "$1" --output=json --home=$SEKAID_HOME | jsonParse
}

# showProposal $(lastProposal) 
function showProposal() {
    sekaid query customgov proposal "$1" --output json --home=$SEKAID_HOME | jsonParse
}

function showProposals() {
    sekaid query customgov proposals --limit=999999999 --output=json --home=$SEKAID_HOME | jsonParse
}

# propAwait $(lastProposal) 
function propAwait() {
    local START_TIME="$(date -u +%s)"
    local ID=""
    local STATUS=""
    local TIMEOUT=""
    if (! $(isNaturalNumber "$1")) ; then
        ID=$(cat) && STATUS=$1 && TIMEOUT=$2
    else
        ID=$1 && STATUS=$2 && TIMEOUT=$3
    fi
    
    local PROP=$(showProposal $ID 2> /dev/null || echo -n "")
    local RESULT=""
    
    if ($(isNullOrEmpty "$PROP")) ; then
        echoErr "ERROR: Proposal $ID was NOT found"
    else
        echoInfo "INFO: Waiting for proposal $ID to be finalized"
        (! $(isNaturalNumber $TIMEOUT)) && TIMEOUT=0
        [[ $TIMEOUT -le 0 ]] && MAX_TIME="∞" || MAX_TIME="$TIMEOUT"
        while : ; do
            local ELAPSED=$(($(date -u +%s) - $START_TIME))
            RESULT=$(showProposal $ID 2> /dev/null | jq ".result" 2> /dev/null | xargs 2> /dev/null || echo -n "")
            ($(isNullOrEmpty "$STATUS")) && ( [ "${RESULT,,}" == "vote_pending" ] || [ "${RESULT,,}" == "vote_result_enactment" ] ) && break
            (! $(isNullOrEmpty "$STATUS")) && [ "${RESULT,,}" == "${STATUS,,}" ] && break
            if [[ $TIMEOUT -gt 0 ]] && [[ $ELAPSED -gt $TIMEOUT ]] ; then
                echoErr "ERROR: Timeout, failed to finalize proposal '$ID' within ${TIMEOUT} s limit"
                return 1
            else
                sleep 0.5
            fi
        done
        echoInfo "INFO: Proposal was finalized ($RESULT)"
    fi
}

# claimValidatorSeat <account> <moniker> <timeout-seconds>
# e.g.: claimValidatorSeat validator "BOB's NODE" 180
function claimValidatorSeat() {
    local ACCOUNT=$1
    local MONIKER=$2
    local TIMEOUT=$3
    ($(isNullOrEmpty $ACCOUNT)) && echoInfo "INFO: Account name was not defined " && return 1
    ($(isNullOrEmpty $MONIKER)) && MONIKER=$(openssl rand -hex 16)
    (! $(isNaturalNumber $TIMEOUT)) && TIMEOUT=180
    sekaid tx customstaking claim-validator-seat --from "$ACCOUNT" --keyring-backend=test --home=$SEKAID_HOME --moniker="$MONIKER" --chain-id=$NETWORK_NAME --broadcast-mode=async --fees=100ukex --yes --output=json | txAwait $TIMEOUT
}

# e.g. showBalance validator
function showBalance() {
    local ADDR=$(showAddress $1)
    local DENOM="$2"
    local RESULT=""
    (! $(isNullOrEmpty $ADDR)) && RESULT=$(sekaid query bank balances "$ADDR" --output=json --home=$SEKAID_HOME 2> /dev/null | jsonParse 2> /dev/null || echo -n "")
    if (! $(isNullOrEmpty $DENOM)) ; then
        RESULT=$(echo $RESULT | showBalance $ADDR | jq '.balances' |  jq ".[] | select(.denom==\"$DENOM\")" | jq '.amount' | xargs 2> /dev/null || echo "0")
        (! $(isNaturalNumber $RESULT)) && RESULT=0
    fi
    echo $RESULT
}

# e.g. sendTokens faucet kiraXXX...XXX 1000 ukex 100 ukex
function sendTokens() {
    local SOURCE=$1
    local DESTINATION=$(showAddress $2)
    local AMOUNT="$3"
    local DENOM="$4"
    local FEE_AMOUNT="$5"
    local FEE_DENOM="$6"

    ($(isNullOrEmpty $FEE_AMOUNT)) && FEE_AMOUNT=100
    ($(isNullOrEmpty $FEE_DENOM)) && FEE_DENOM="ukex"

    echoInfo "INFO: Sending $AMOUNT $DENOM | $SOURCE -> $DESTINATION"
    OLD_BALANCE_SRC=$(showBalance "$SOURCE" "$DENOM") && (! $(isNaturalNumber $OLD_BALANCE_SRC)) && OLD_BALANCE_SRC=0
    OLD_BALANCE_SRC_FEE=$(showBalance "$SOURCE" "$FEE_DENOM") && (! $(isNaturalNumber $OLD_BALANCE_SRC_FEE)) && OLD_BALANCE_SRC_FEE=0
    OLD_BALANCE_DEST=$(showBalance "$DESTINATION" "$DENOM") && (! $(isNaturalNumber $OLD_BALANCE_DEST)) && OLD_BALANCE_DEST=0

    sekaid tx bank send $SOURCE $DESTINATION "${AMOUNT}${DENOM}" --keyring-backend=test --chain-id=$NETWORK_NAME --fees "${FEE_AMOUNT}${FEE_DENOM}" --output=json --yes --home=$SEKAID_HOME | txAwait 180

    NEW_BALANCE_SRC=$(showBalance "$SOURCE" "$DENOM") && (! $(isNaturalNumber $NEW_BALANCE_SRC)) && NEW_BALANCE_SRC=0
    NEW_BALANCE_SRC_FEE=$(showBalance "$SOURCE" "$FEE_DENOM") && (! $(isNaturalNumber $NEW_BALANCE_SRC_FEE)) && NEW_BALANCE_SRC_FEE=0
    NEW_BALANCE_DEST=$(showBalance "$DESTINATION" "$DENOM") && (! $(isNaturalNumber $NEW_BALANCE_DEST)) && NEW_BALANCE_DEST=0

    [ "$OLD_BALANCE_SRC" != "$NEW_BALANCE_SRC" ] && echoInfo "INFO:  SRC. balance change $(showAddress $1) | $OLD_BALANCE_SRC $DENOM -> $NEW_BALANCE_SRC $DENOM"
    [ "$OLD_BALANCE_SRC_FEE" != "$NEW_BALANCE_SRC_FEE" ] && [ "$DENOM" != "$FEE_DENOM" ] && \
    echoInfo "INFO:  SRC. balance change $(showAddress $1) | $OLD_BALANCE_SRC_FEE $FEE_DENOM -> $NEW_BALANCE_SRC_FEE $FEE_DENOM"
    [ "$OLD_BALANCE_DEST" != "$NEW_BALANCE_DEST" ] && echoInfo "INFO: DEST. balance change $DESTINATION | $OLD_BALANCE_DEST $DENOM -> $NEW_BALANCE_DEST $DENOM"
}

# e.g. showStatus -> { ... }
function showStatus() {
    echo $(sekaid status 2>&1 | bash-utils jsonParse "" 2>/dev/null || echo -n "")
}

# e.g. showBlockHeight -> 123
function showBlockHeight() {
    SH_LATEST_BLOCK_HEIGHT=$(showStatus | bash-utils jsonParse "SyncInfo.latest_block_height" 2>/dev/null || echo -n "")
    ($(bash-utils isNaturalNumber "$SH_LATEST_BLOCK_HEIGHT")) && echo $SH_LATEST_BLOCK_HEIGHT || echo ""
}

# awaitBlocks <number-of-blocks>
# e.g. awaitBlocks 5
function awaitBlocks() {
    local BLOCKS=$1
    (! $(bash-utils isNaturalNumber $BLOCKS)) && bash-utils echoErr "ERROR: Number of blocks to await was NOT defined" && return 1
    local SH_START_BLOCK=""
    while : ; do
        local SH_NEW_BLOCK=$(showBlockHeight)
        (! $(bash-utils isNaturalNumber $SH_NEW_BLOCK)) && sleep 0.5 && continue
        (! $(bash-utils isNaturalNumber "$SH_START_BLOCK")) && SH_START_BLOCK=$SH_NEW_BLOCK
        local SH_DELTA=$(($SH_NEW_BLOCK - $SH_START_BLOCK))
        [ $SH_DELTA -ge $BLOCKS ] && break
        sleep 0.5
    done
}

# e.g. showCatchingUp -> false
function showCatchingUp() {
    local SH_CATCHING_UP=$(showStatus | jsonParse "SyncInfo.catching_up" 2>/dev/null || echo -n "")
    ($(isBoolean "$SH_CATCHING_UP")) && echo "${SH_CATCHING_UP,,}" || echo ""
}

# activateValidator <account> <timeout-seconds>
# e.g. activateValidator validator 180
function activateValidator() {
    local ACCOUNT=$1
    local TIMEOUT=$2
    (! $(isNaturalNumber $TIMEOUT)) && TIMEOUT=180
    ($(isNullOrEmpty $ACCOUNT)) && echoInfo "INFO: Account name was not defined " && return 1
    sekaid tx customslashing activate --from "$ACCOUNT" --chain-id=$NETWORK_NAME --keyring-backend=test --home=$SEKAID_HOME --fees 1000ukex --yes --broadcast-mode=async --log_format=json --output=json | txAwait $TIMEOUT
}

# pauseValidator <account> <timeout-seconds>
# e.g. pauseValidator validator 180
function pauseValidator() {
    local ACCOUNT=$1
    local TIMEOUT=$2
    (! $(isNaturalNumber $TIMEOUT)) && TIMEOUT=180
    ($(isNullOrEmpty $ACCOUNT)) && echoInfo "INFO: Account name was not defined " && return 1
    sekaid tx customslashing pause --from "$ACCOUNT" --chain-id=$NETWORK_NAME --keyring-backend=test --home=$SEKAID_HOME --fees 100ukex --yes --broadcast-mode=async --log_format=json --output=json | txAwait $TIMEOUT
}

# unpauseValidator <account> <timeout-seconds>
# e.g. unpauseValidator validator 180
function unpauseValidator() {
    local ACCOUNT=$1
    local TIMEOUT=$2
    (! $(isNaturalNumber $TIMEOUT)) && TIMEOUT=180
    ($(isNullOrEmpty $ACCOUNT)) && echoInfo "INFO: Account name was not defined " && return 1
    sekaid tx customslashing unpause --from "$ACCOUNT" --chain-id=$NETWORK_NAME --keyring-backend=test --home=$SEKAID_HOME --fees 100ukex --yes --broadcast-mode=async --log_format=json --output=json | txAwait $TIMEOUT
}

# clearPermission <account> <permission> <address> <timeout-seconds>
# e.g. clearPermission validator 11 kiraXXX..YYY 180
function clearPermission() {
    local ACCOUNT=$1
    local PERM=$2
    local ADDR=$(showAddress $3)
    local TIMEOUT=$4
    ($(isNullOrEmpty $ACCOUNT)) && echoInfo "INFO: Account name was not defined '$1'" && return 1
    ($(isNullOrEmpty $ADDR)) && echoInfo "INFO: Address name was not defined '$3'" && return 1
    (! $(isNaturalNumber $PERM)) && echoInfo "INFO: Invalid permission id '$PERM' " && return 1
    (! $(isNaturalNumber $TIMEOUT)) && TIMEOUT=180
    if ($(isPermBlacklisted $ADDR $PERM)) ; then
        echoInfo "INFO: Permission '$PERM' is blacklisted and will be removed from the blacklist, please wait..."
        sekaid tx customgov permission remove-blacklisted --from "$ACCOUNT" --keyring-backend=test --permission="$PERM" --addr="$ADDR" --chain-id=$NETWORK_NAME --home=$SEKAID_HOME --fees=100ukex --yes --broadcast-mode=async --log_format=json --output=json | txAwait $TIMEOUT
    elif ($(isPermWhitelisted $ADDR $PERM)) ; then
        echoInfo "INFO: Permission '$PERM' is whitelisted and will be removed from the whitelist, please wait..."
        sekaid tx customgov permission remove-whitelisted --from "$ACCOUNT" --keyring-backend=test --permission="$PERM" --addr="$ADDR" --chain-id=$NETWORK_NAME --home=$SEKAID_HOME --fees=100ukex --yes --broadcast-mode=async --log_format=json --output=json | txAwait $TIMEOUT
    else
        echoInfo "INFO: Permission '$PERM' was never present or already cleared."
    fi
}

# whitelistPermission <account> <permission> <address> <timeout-seconds>
# e.g. whitelistPermission validator 11 kiraXXX..YYY 180
function whitelistPermission() {
    local KM_ACC=$1
    local PERM=$2
    local ADDR=$(showAddress $3)
    local TIMEOUT=$4
    ($(isNullOrEmpty $KM_ACC)) && echoInfo "INFO: Account name was not defined '$1'" && return 1
    ($(isNullOrEmpty $ADDR)) && echoInfo "INFO: Address name was not defined '$3'" && return 1
    (! $(isNaturalNumber $PERM)) && echoInfo "INFO: Invalid permission id '$PERM' " && return 1
    (! $(isNaturalNumber $TIMEOUT)) && TIMEOUT=180
    if ($(isPermWhitelisted $ADDR $PERM)) ; then
        echoWarn "WARNING: Address '$ADDR' already has assigned permission '$PERM'"
    else
        if ($(isPermBlacklisted $ADDR $PERM)) ; then
            echoWarn "WARNING: Address '$ADDR' has blacklisted permission '$PERM', attempting to clear..."
            clearPermission $KM_ACC $PERM $ADDR $TIMEOUT
        fi

        sekaid tx customgov permission whitelist --from "$KM_ACC" --keyring-backend=test --permission="$PERM" --addr="$ADDR" --chain-id=$NETWORK_NAME --home=$SEKAID_HOME --fees=100ukex --yes --broadcast-mode=async --log_format=json --output=json | txAwait $TIMEOUT
    fi
}

# blacklisttPermission <account> <permission> <address> <timeout-seconds>
# e.g. blacklisttPermission validator 11 kiraXXX..YYY 180
function blacklistPermission() {
    local KM_ACC=$1
    local PERM=$2
    local ADDR=$(showAddress $3)
    local TIMEOUT=$4
    ($(isNullOrEmpty $KM_ACC)) && echoInfo "INFO: Account name was not defined '$1'" && return 1
    ($(isNullOrEmpty $ADDR)) && echoInfo "INFO: Address name was not defined '$3'" && return 1
    (! $(isNaturalNumber $PERM)) && echoInfo "INFO: Invalid permission id '$PERM' " && return 1
    (! $(isNaturalNumber $TIMEOUT)) && TIMEOUT=180
    if ($(isPermBlacklisted $ADDR $PERM)) ; then
        echoWarn "WARNING: Address '$ADDR' already has blacklisted permission '$PERM'"
    else
        if ($(isPermWhitelisted $ADDR $PERM)) ; then
            echoWarn "WARNING: Address '$ADDR' has whitelisted permission '$PERM', attempting to clear..."
            clearPermission $KM_ACC $PERM $ADDR $TIMEOUT
        fi

        sekaid tx customgov permission blacklist --from "$KM_ACC" --keyring-backend=test --permission="$PERM" --addr="$ADDR" --chain-id=$NETWORK_NAME --home=$SEKAID_HOME --fees=100ukex --yes --broadcast-mode=async --log_format=json --output=json | txAwait $TIMEOUT
    fi
}

function showCurrentPlan() {
    sekaid query upgrade current-plan --output=json --chain-id=$NETWORK_NAME --home=$SEKAID_HOME
}

function showNextPlan() {
    sekaid query upgrade next-plan --output=json --chain-id=$NETWORK_NAME --home=$SEKAID_HOME
}

# showIdentityRecord <account> <key> // shows all or a single key
# e.g. showIdentityRecord validator "mykey"
# e.g. showIdentityRecord validator 15
function showIdentityRecord() {
    local KM_ACC=$(showAddress $1)
    local KM_KEY=$2 && [ "$KM_KEY" == "*" ] && KM_KEY=""

    ($(isNullOrEmpty $KM_ACC)) && echoErr "ERROR: Account name or address '$1' is invalid" && return 1
    if ($(isNullOrEmpty $KM_KEY)) ; then
        sekaid query customgov identity-records-by-addr $KM_ACC --output=json --home=$SEKAID_HOME | jq 2> /dev/null || echo ""
    else
        if ($(isNumber $KM_KEY)) ; then
            sekaid query customgov identity-records-by-addr $KM_ACC --output=json --home=$SEKAID_HOME | jq ".records | .[] | select(.id==\"$KM_KEY\")" 2> /dev/null || echo ""
        else
            sekaid query customgov identity-records-by-addr $KM_ACC --output=json --home=$SEKAID_HOME | jq ".records | .[] | select(.key==\"$KM_KEY\")" 2> /dev/null || echo ""
        fi
    fi
}

# upsertIdentityRecord <account> <key> <value> <timeout-seconds>
# e.g. upsertIdentityRecord validator "mykey" "My Value" 180
function upsertIdentityRecord() {
    local KM_ACC=$(showAddress $1)
    local IR_KEY=$2
    local IR_VAL=$3
    local TIMEOUT=$4
    ($(isNullOrEmpty $KM_ACC)) && echoErr "ERROR: Account name was NOT defined " && return 1
    ($(isNullOrEmpty $IR_KEY)) && echoErr "ERROR: Key was NOT defined " && return 1
    (! $(isNaturalNumber $TIMEOUT)) && TIMEOUT=180

    if ($(isNullOrEmpty $IR_VAL)) ; then
        sekaid tx customgov delete-identity-records --keys="$IR_KEY" --from=$KM_ACC --keyring-backend=test --home=$SEKAID_HOME --chain-id=$NETWORK_NAME --fees=100ukex --yes --broadcast-mode=async --log_format=json --output=json | txAwait $TIMEOUT
    else
        sekaid tx customgov register-identity-records --infos-json="{\"$IR_KEY\":\"$IR_VAL\"}" --from=$KM_ACC --keyring-backend=test --home=$SEKAID_HOME --chain-id=$NETWORK_NAME --fees=100ukex --yes --broadcast-mode=async --log_format=json --output=json | txAwait $TIMEOUT
    fi
}

# verifyIdentityRecord <account> <verifier-address> <one-or-many-comma-separated-keys/ids> <tip> <timeout>
# e.g. verifyIdentityRecord validator $(showAddress test) "mykey,mykey2" "200ukex" 180
function verifyIdentityRecord() {
    local KM_ACC=$1
    local KM_VER=$(showAddress $2)
    local KM_KEYS=$3
    local KM_TIP=$4
    local TIMEOUT=$5
    ($(isNullOrEmpty $KM_ACC)) && echoErr "ERROR: Account name was NOT defined '$1'" && return 1
    ($(isNullOrEmpty $KM_VER)) && echoErr "ERROR: Verifier address '$2' is invalid" && return 1
    ($(isNullOrEmpty $KM_KEYS)) && echoErr "ERROR: Record keys to verify were NOT specified '$3'" && return 1
    (! $(isNaturalNumber $TIMEOUT)) && TIMEOUT=180

    local FINAL_IR_KEYS=""
    for irkey in $(echo $KM_KEYS | sed "s/,/ /g") ; do
        local irkey_id=$(showIdentityRecord $KM_ACC "${irkey,,}" | jsonParse ".id" 2> /dev/null || echo "")
        ($(isNullOrEmpty $irkey_id)) && echoErr "ERROR: Key '$irkey' is invalid or was NOT found" && return 1
        [ ! -z "$FINAL_IR_KEYS" ] && FINAL_IR_KEYS="${FINAL_IR_KEYS},"
        FINAL_IR_KEYS="${FINAL_IR_KEYS}${irkey_id}"
    done
    ($(isNullOrEmpty $FINAL_IR_KEYS)) && echoErr "ERROR: No valid record keys were found" && return 1

    echoInfo "INFO: Sending request to verify '$FINAL_IR_KEYS'"
    sekaid tx customgov request-identity-record-verify --verifier="$KM_VER" --record-ids="$FINAL_IR_KEYS" --from=$KM_ACC --tip="$KM_TIP" --keyring-backend=test --home=$SEKAID_HOME --chain-id=$NETWORK_NAME --fees=100ukex --yes --broadcast-mode=async --log_format=json --output=json | txAwait $TIMEOUT
}

# showIdentityVerificationRequests <verifier-account/address> <requester-address>
# e.g. showIdentityVerificationRequests validator $(showAddress test)
function showIdentityVerificationRequests() {
    local KM_ACC=$(showAddress $1)
    local KM_REQ=$(showAddress $2)

    ($(isNullOrEmpty $KM_ACC)) && echoErr "ERROR: Account name or address '$1' is invalid" && return 1
    if ($(isNullOrEmpty $KM_REQ)) ; then
        sekaid query customgov identity-record-verify-requests-by-approver $KM_ACC --output=json --home=$SEKAID_HOME  | jq 2> /dev/null || echo ""
    else
        sekaid query customgov identity-record-verify-requests-by-approver $KM_ACC --output=json --home=$SEKAID_HOME  | jq ".verify_records | .[] | select(.address==\"$KM_REQ\")" 2> /dev/null || echo ""
    fi
}

# approveIdentityVerificationRequest <account> <id> <timeout>
# e.g. approveIdentityVerificationRequest validator 1 180
function approveIdentityVerificationRequest() {
    local ACCOUNT=$1
    local KM_REQ=$2
    local TIMEOUT=$3 && (! $(isNaturalNumber $TIMEOUT)) && TIMEOUT=180
    ($(isNullOrEmpty $ACCOUNT)) && echoErr "ERROR: Account name was NOT defined '$1'" && return 1
    (! $(isNaturalNumber $KM_REQ)) && echoErr "ERROR: Request Id must be a valid natural number, but got '$KM_REQ'" && return 1
    
    sekaid tx customgov handle-identity-records-verify-request $KM_REQ --approve="true" --from=$ACCOUNT --keyring-backend=test --home=$SEKAID_HOME --chain-id=$NETWORK_NAME --fees=100ukex --yes --broadcast-mode=async --log_format=json --output=json | txAwait $TIMEOUT
}

# rejectIdentityVerificationRequest <account> <id> <timeout>
# e.g. rejectIdentityVerificationRequest validator 1 180
function rejectIdentityVerificationRequest() {
    local ACCOUNT=$1
    local KM_REQ=$2
    local TIMEOUT=$3 && (! $(isNaturalNumber $TIMEOUT)) && TIMEOUT=180
    ($(isNullOrEmpty $ACCOUNT)) && echoErr "ERROR: Account name was NOT defined " && return 1
    (! $(isNaturalNumber $KM_REQ)) && echoErr "ERROR: Request Id must be a valid natural number, but got '$KM_REQ'" && return 1
    
    sekaid tx customgov handle-identity-records-verify-request $KM_REQ --approve="false" --from=$ACCOUNT --keyring-backend=test --home=$SEKAID_HOME --chain-id=$NETWORK_NAME --fees=100ukex --yes --broadcast-mode=async --log_format=json --output=json | txAwait $TIMEOUT
}

# upsertDataRegistry <account> <key> <value> <file-type>
# e.g: upsertDataRegistry validator "code_of_conduct" "https://raw.githubusercontent.com/KiraCore/sekai/master/sekai-env.sh" "text"
function upsertDataRegistry() {
    local ACCOUNT=$1
    local KEY=$2
    local VALUE=$3
    local FILETYPE=$4

    local DOWNLOAD_SUCCESS="true"
    local TMP_FILE="/tmp/data-registry.tmp"
    rm -fv $TMP_FILE
    wget "$VALUE" -O $TMP_FILE || DOWNLOAD_SUCCESS="false"

    if [ "${DOWNLOAD_SUCCESS,,}" != "true" ] ; then
        echoErr "ERROR: Resource '$VALUE' was NOT found, failed to create proposal"
        return 1
    else
        echoInfo "SUCCESS: Resource '$VALUE' was found."
        local CHECKSUM=$(sha256 $TMP_FILE)
        local SIZE=$(fileSize $TMP_FILE)
        echoInfo "INFO: Voting YES on proposal $PROPOSAL with account $ACCOUNT"
        sekaid tx customgov proposal upsert-data-registry "$KEY" "$CHECKSUM" "$VALUE" "$FILETYPE" "$SIZE" --title="Upserting Data Registry key '$KEY'" --description="Assign value '$VALUE' to key '$KEY'" --from=$ACCOUNT --chain-id=$NETWORK_NAME --keyring-backend=test  --fees=100ukex --yes --log_format=json --broadcast-mode=async --output=json --home=$SEKAID_HOME  | txAwait
    fi
}

function showDataRegistryKeys() {
    sekaid query customgov data-registry-keys --page-key 1000000 --output=json | jq ".keys"
}

function showDataRegistryKey() {
    sekaid query customgov data-registry "$1" --output=json 2> /dev/null || echo ""
}

# setNetworkProperty <key> <value>
# e.g: setNetworkProperty validator "MIN_TX_FEE" "99"
function setNetworkProperty() {
    local ACCOUNT=$1
    local KEY=$2
    local VALUE=$3

    sekaid tx customgov proposal set-network-property "${KEY^^}" "$VALUE" --title="Upserting Network Property '$KEY'" --description="Assign value '$VALUE' to property '$KEY'" --from=$ACCOUNT --chain-id=$NETWORK_NAME --keyring-backend=test  --fees=100ukex --yes --log_format=json --broadcast-mode=async --output=json --home=$SEKAID_HOME  | txAwait
}

# createRole <account> <role-name> <role-description>
# e.g. createRole validator validator "Role enabling to claim validator seat and perform essential gov. functions"
function createRole() {
    local ACCOUNT=$1
    local NAME=$2
    local DESCRIPTION=$3
    local TIMEOUT=$4
    ($(isNullOrEmpty $ACCOUNT)) && echoInfo "INFO: Account name was not defined '$1'" && return 1
    ($(isNullOrEmpty $NAME)) && echoInfo "INFO: Invalid role name '$NAME' " && return 1
    (! $(isNaturalNumber $TIMEOUT)) && TIMEOUT=180
    sekaid tx customgov role create "$NAME" "$DESCRIPTION" --from "$ACCOUNT" --keyring-backend=test --chain-id=$NETWORK_NAME --home=$SEKAID_HOME --fees=100ukex --yes --broadcast-mode=async --log_format=json --output=json | txAwait $TIMEOUT
}

# showRoles <account-or-address>
# e.g. showRoles validator 
function showRoles() {
    local ADDRESS=$(showAddress $1)
    if ($(isNullOrEmpty $ADDRESS)) ; then
        echo $(sekaid query customgov all-roles --output=json --home=$SEKAID_HOME 2> /dev/null | jsonParse 2> /dev/null || echo -n "") && echo -n ""
    else
        echo $(sekaid query customgov roles $ADDRESS --output=json --home=$SEKAID_HOME 2> /dev/null | jsonParse 2> /dev/null || echo -n "") && echo -n ""
    fi
}

function showRole() {
    echo $(sekaid query customgov role "$1" --output=json --home=$SEKAID_HOME 2> /dev/null | jsonParse 2> /dev/null || echo -n "") && echo -n ""
}

# setProposalsDurations <account> <comma-separated-proposals> <comma-separated-time-values>
# e.g: setProposalsDurations validator "UpsertDataRegistry,SetNetworkProperty" "300,300"
function setProposalsDurations() {
    local ACCOUNT=$1
    local PROPOSALS=$2
    local DURATIONS=$3
    ($(isNullOrEmpty $PROPOSALS)) && echoInfo "INFO: Proposals were NOT defined '$2'" && return 1
    ($(isNullOrEmpty $DURATIONS)) && echoInfo "INFO: Durations were NOT defined '$3'" && return 1

    sekaid tx customgov proposal set-proposal-durations-proposal "$PROPOSALS" "$DURATIONS" --title="Update proposals duration " --description="Set durations of '[$PROPOSALS]' to '[$DURATIONS]' seconds" --from "$ACCOUNT" --chain-id=$NETWORK_NAME --keyring-backend=test  --fees=100ukex --yes --log_format=json --broadcast-mode=async --output=json --home=$SEKAID_HOME | txAwait
}

function showProposalsDurations() {
    echo $(sekaid query customgov all-proposal-durations --output=json --home=$SEKAID_HOME 2> /dev/null | jsonParse 2> /dev/null || echo -n "") && echo -n ""
}

function showPoorNetworkMessages() {
    echo $(sekaid query customgov poor-network-messages --output=json --home=$SEKAID_HOME 2> /dev/null | jsonParse 2> /dev/null || echo -n "") && echo -n ""
}

# setPoorNetworkMessages <account> <comma-transaction-types>
# e.g: setPoorNetworkMessages validator "submit_evidence,submit-proposal,vote-proposal,claim-councilor,set-network-properties,claim-validator,activate,pause,unpause"
function setPoorNetworkMessages() {
    local ACCOUNT=$1
    local MESSAGES=$2
    ($(isNullOrEmpty $ACCOUNT)) && echoInfo "INFO: Account was NOT defined '$1'" && return 1
    ($(isNullOrEmpty $MESSAGES)) && echoInfo "INFO: Allowed network messages were NOT defined '$2'" && return 1

    sekaid tx customgov proposal set-poor-network-msgs "$MESSAGES" --title="Update poor network messages" --description="Allowing submission of '[$MESSAGES]' during poor network conditions" --from "$ACCOUNT" --chain-id=$NETWORK_NAME --keyring-backend=test  --fees=100ukex --yes --log_format=json --broadcast-mode=async --output=json --home=$SEKAID_HOME | txAwait
}

# showExecutionFee <transaction-type>
# e.g.: showExecutionFee <transaction-type>
function showExecutionFee() {
    local TRANSACTION_TYPE=$1
    local MIN_TX_FEE=$(showNetworkProperties | jsonParse "properties.min_tx_fee")
    local EXECUTION_FEE=$(sekaid query customgov execution-fee "$TRANSACTION_TYPE" --output=json --home=$SEKAID_HOME 2> /dev/null | jsonParse "fee.execution_fee" 2> /dev/null || echo -n "")
    (! $(isNaturalNumber $EXECUTION_FEE)) && EXECUTION_FEE=$MIN_TX_FEE
    [ $MIN_TX_FEE -gt $EXECUTION_FEE ] && echo "$MIN_TX_FEE" || echo "$EXECUTION_FEE" 
}

# setExecutionFee <account> <tx-type> <execution-fee> <failure-fee> <tx-timeout>
# e.g.: setExecutionFee validator pause 100 200 60
function setExecutionFee() {
    local ACCOUNT=$1
    local TX_TYPE=$2
    local EXECUTION_FEE=$3
    local FAILURE_FEE=$4
    local TX_TIMEOUT=$5
    ($(isNullOrEmpty $ACCOUNT)) && echoInfo "INFO: Account was NOT defined '$1'" && return 1
    ($(isNullOrEmpty $TX_TYPE)) && echoInfo "INFO: Transaction type was NOT defined '$2'" && return 1
    (! $(isNaturalNumber $EXECUTION_FEE)) && echoInfo "INFO: Invalid execution fee amount '$3'" && return 1
    (! $(isNaturalNumber $FAILURE_FEE)) && echoInfo "INFO: Invalid failure fee amount '$4'" && return 1
    (! $(isNaturalNumber $TX_TIMEOUT)) && echoInfo "INFO: Invalid tx timeout '$5'" && return 1
    local TX_FEE="$(showExecutionFee 'set-execution-fee')ukex"
    sekaid tx customgov set-execution-fee --transaction_type="$TX_TYPE" --execution_fee="$EXECUTION_FEE" --failure_fee="$FAILURE_FEE" --timeout="$TX_TIMEOUT" --from "$ACCOUNT" --chain-id=$NETWORK_NAME --keyring-backend=test  --fees=$TX_FEE --yes --log_format=json --broadcast-mode=async --output=json --home=$SEKAID_HOME | txAwait
}

# resetRanks <account>
# e.g: resetRanks validator
function resetRanks() {
    local ACCOUNT=$1
    ($(isNullOrEmpty $ACCOUNT)) && echoInfo "INFO: Account was NOT defined '$1'" && return 1

    sekaid tx customslashing proposal-reset-whole-validator-rank --title="Ranks reset" --description="Reseting ranks or all validator nodes" --from "$ACCOUNT" --chain-id=$NETWORK_NAME --keyring-backend=test  --fees=100ukex --yes --log_format=json --broadcast-mode=async --output=json --home=$SEKAID_HOME | txAwait
}

# showValidator <account/kira-address/val-address>
# e.g. showValidator validator
# e.g. showValidator kiraXXXXXXXXXXX
# e.g. showValidator kiravaloperXXXXXXXXXXX
function showValidator() {
    local VAL_ADDR="${1,,}"
    if [[ $VAL_ADDR == kiravaloper* ]] ; then
        VAL_STATUS=$(sekaid query customstaking validator --val-addr="$VAL_ADDR" --output=json --home=$SEKAID_HOME 2> /dev/null | jsonParse 2> /dev/null || echo -n "")
    else
        local ADDRESS=$(showAddress $VAL_ADDR)
        VAL_STATUS=$(sekaid query customstaking validator --addr="$ADDRESS" --output=json --home=$SEKAID_HOME 2> /dev/null | jsonParse 2> /dev/null || echo -n "") 
    fi
    echo $VAL_STATUS
}

function showTokenAliases() {
    echo $(sekaid query tokens all-aliases --output=json --home=$SEKAID_HOME 2> /dev/null | jsonParse 2> /dev/null || echo -n "") 
}

function showTokenRates() {
    echo $(sekaid query tokens all-rates --output=json --home=$SEKAID_HOME 2> /dev/null | jsonParse 2> /dev/null || echo -n "") 
}

# setTokenRate <account> <denom> <rate> <is-fee-token>
function setTokenRate() {
    local ACCOUNT=$1
    local DENOM=$2
    local RATE=$3
    local FEE_PAYMENT=$4
    ($(isNullOrEmpty $ACCOUNT)) && echoInfo "INFO: Account was NOT defined '$1'" && return 1
    ($(isNullOrEmpty $DENOM)) && echoInfo "INFO: Token Denom was NOT defined '$2'" && return 1
    ($(isNullOrEmpty $RATE)) && echoInfo "INFO: Token Exchange Rate was NOT defined '$3'" && return 1
    (! $(isBoolean $FEE_PAYMENT)) && echoInfo "INFO: It must be indicated if token is or is NOT a payment method, but got '$4'" && return 1

    sekaid tx tokens proposal-upsert-rate --denom="$DENOM" --rate="$RATE" --fee_payments="$FEE_PAYMENT" --title="Set exchange rate of '$DENOM'" --description="Fee payments will be set at the rate of $RATE $DENOM == 1 KEX. Set '$FEE_PAYMENT' to indicate if $DENOM is a payment method." --from "$ACCOUNT" --chain-id=$NETWORK_NAME --keyring-backend=test  --fees=100ukex --yes --log_format=json --broadcast-mode=async --output=json --home=$SEKAID_HOME | txAwait
}

function setTokensBlackWhiteList() {
    local ACCOUNT=$1
    local IS_BLACKLIST=$2
    local IS_ADD=$3
    local TOKENS=$4
    ($(isNullOrEmpty $ACCOUNT)) && echoInfo "INFO: Account was NOT defined '$1'" && return 1
    (! $(isBoolean $IS_BLACKLIST)) && echoInfo "INFO: Is Blacklist parameter must be a boolean, but got '$2'" && return 1
    (! $(isBoolean $IS_ADD)) && echoInfo "INFO: Is Add parameter must be a boolean, but got '$3'" && return 1
    ($(isNullOrEmpty $TOKENS)) && echoInfo "INFO: Tokens to add/remove from black/white list were NOT defined '$4'" && return 1

    sekaid tx tokens proposal-update-tokens-blackwhite --is_add="$IS_ADD" --is_blacklist="$IS_BLACKLIST" --tokens="$TOKENS" --title="Update Tokens Black/White-list" --description="Is Blacklist: '$IS_BLACKLIST', Is Add: $IS_ADD, Tokens: '$TOKENS'" --from "$ACCOUNT" --chain-id=$NETWORK_NAME --keyring-backend=test  --fees=100ukex --yes --log_format=json --broadcast-mode=async --output=json --home=$SEKAID_HOME | txAwait
}

# whitelistAddTokenTransfers <account> <tokens>
function transfersWhitelistAddTokens() {
    setTokensBlackWhiteList "$1" "false" "true" "$2"
}

function transfersWhitelistRemoveTokens() {
    setTokensBlackWhiteList "$1" "false" "false" "$2"
}

function transfersBlacklistAddTokens() {
    setTokensBlackWhiteList "$1" "true" "true" "$2"
}

function transfersBlacklistRemoveTokens() {
    setTokensBlackWhiteList "$1" "true" "false" "$2"
}

function showTokenTransferBlackWhiteList() {
    echo $(sekaid query tokens token-black-whites --output=json --home=$SEKAID_HOME 2> /dev/null | jsonParse 2> /dev/null || echo -n "") 
}

function unjail() {
    local ACCOUNT=$1
    local ADDRESS=$2
    local REFERENCE=$3
    ADDRESS=$(showValidator "$ADDRESS" | jsonParse "val_key" 2> /dev/null || echo -n "");

    ($(isNullOrEmpty $ACCOUNT)) && echoInfo "INFO: Account was NOT defined '$1'" && return 1
    ($(isNullOrEmpty $ADDRESS)) && echoInfo "INFO: Validator Address to unjail was NOT defined or could NOT be found '$2'" && return 1
    # ($(isNullOrEmpty $REFERENCE)) && echoInfo "INFO: Unjail reference should NOT be empty '$3'" && return 1
    sekaid tx customstaking proposal unjail-validator "$ADDRESS" "$REFERENCE" --title="Unjail validator '$ADDRESS'" --description="Proposal to unjail '$ADDRESS' due to his unintentional fault" --from "$ACCOUNT" --chain-id=$NETWORK_NAME --keyring-backend=test  --fees=100ukex --yes --log_format=json --broadcast-mode=async --output=json --home=$SEKAID_HOME | txAwait
}

# showPermissions validator
function showRolePermissions() {
    echo $(sekaid query customgov role "$1" --output=json --home=$SEKAID_HOME 2> /dev/null | jsonParse "permissions" 2> /dev/null || echo -n "") && echo -n ""
}

function isRolePermBlacklisted() {
    local ROLE=$1
    local PERM=$2
    if (! $(isNaturalNumber $PERM)) || ($(isNullOrEmpty $ROLE)) ; then
        echo "false"
    else
        INDEX=$(showRolePermissions $ROLE 2> /dev/null | jq ".blacklist | index($PERM)" 2> /dev/null || echo -n "")
        ($(isNaturalNumber $INDEX)) && echo "true" || echo "false"
    fi
}

function isRolePermWhitelisted() {
    local ROLE=$1
    local PERM=$2
    if (! $(isNaturalNumber $PERM)) || ($(isNullOrEmpty $ROLE)) ; then
        echo "false"
    else
        INDEX=$(showRolePermissions $ROLE 2> /dev/null | jq ".whitelist | index($PERM)" 2> /dev/null || echo -n "")
        if ($(isNaturalNumber $INDEX)) && (! $(isRolePermBlacklisted $ROLE $PERM)) ; then
            echo "true" 
        else
            echo "false"
        fi
    fi
}

# roleClearPermission <account> <role> <permission> <timeout-seconds>
# e.g. roleClearPermission validator validator 2 180
function roleClearPermission() {
    local ACCOUNT=$1
    local ROLE=$2
    local PERM=$3
    local TIMEOUT=$4
    ($(isNullOrEmpty $ACCOUNT)) && echoInfo "INFO: Account name was not defined '$ACCOUNT'" && return 1
    ($(isNullOrEmpty $ROLE)) && echoInfo "INFO: Role name was not defined '$ROLE'" && return 1
    (! $(isNaturalNumber $PERM)) && echoInfo "INFO: Invalid permission id '$PERM' " && return 1
    (! $(isNaturalNumber $TIMEOUT)) && TIMEOUT=180
    if ($(isRolePermBlacklisted "$ROLE" "$PERM")) ; then
        echoInfo "INFO: Permission '$PERM' is blacklisted and will be removed from the '$ROLE' role blacklist, please wait..."
        sekaid tx customgov role remove-blacklisted-permission "$ROLE" "$PERM" --from "$ACCOUNT" --keyring-backend=test --chain-id=$NETWORK_NAME --home=$SEKAID_HOME --fees=100ukex --yes --broadcast-mode=async --log_format=json --output=json | txAwait $TIMEOUT
    elif ($(isRolePermWhitelisted "$ROLE" "$PERM")) ; then
        echoInfo "INFO: Permission '$PERM' is whitelisted and will be removed from the '$ROLE' role whitelist, please wait..."
        sekaid tx customgov role remove-whitelisted-permission "$ROLE" "$PERM" --from "$ACCOUNT" --keyring-backend=test --chain-id=$NETWORK_NAME --home=$SEKAID_HOME --fees=100ukex --yes --broadcast-mode=async --log_format=json --output=json | txAwait $TIMEOUT
    else
        echoInfo "INFO: Permission '$PERM' was never present in the '$ROLE' role blacklist/whitelist or was already cleared."
    fi
}

## roleWhitelistPermission <account> <role> <permission> <timeout-seconds>
## e.g. roleWhitelistPermission validator validator 29 kiraXXX..YYY 180
function roleWhitelistPermission() {
    local ACCOUNT=$1
    local ROLE=$2
    local PERM=$3
    local TIMEOUT=$4
    ($(isNullOrEmpty $ACCOUNT)) && echoInfo "INFO: Account name was not defined '$1'" && return 1
    ($(isNullOrEmpty $ROLE)) && echoInfo "INFO: Role name was not defined '$ROLE'" && return 1
    (! $(isNaturalNumber $PERM)) && echoInfo "INFO: Invalid permission id '$PERM' " && return 1
    (! $(isNaturalNumber $TIMEOUT)) && TIMEOUT=180
    if ($(isRolePermWhitelisted "$ROLE" "$PERM")) ; then
        echoWarn "WARNING: Role '$ROLE' whitelist already has assigned permission '$PERM'"
    else
        if ($(isRolePermBlacklisted "$ROLE" "$PERM")) ; then
            echoWarn "WARNING: Role '$ROLE' has blacklisted permission '$PERM', attempting to clear..."
            roleClearPermission $ACCOUNT $ROLE $PERM $TIMEOUT
        fi

        sekaid tx customgov role whitelist-permission "$ROLE" "$PERM" --from "$ACCOUNT" --keyring-backend=test --chain-id=$NETWORK_NAME --home=$SEKAID_HOME --fees=100ukex --yes --broadcast-mode=async --log_format=json --output=json | txAwait $TIMEOUT
    fi
}

## roleBlacklistPermission <account> <role> <permission> <timeout-seconds>
## e.g. roleBlacklistPermission validator validator 29 kiraXXX..YYY 180
function roleBlacklistPermission() {
    local ACCOUNT=$1
    local ROLE=$2
    local PERM=$3
    local TIMEOUT=$4
    ($(isNullOrEmpty $ACCOUNT)) && echoInfo "INFO: Account name was not defined '$1'" && return 1
    ($(isNullOrEmpty $ROLE)) && echoInfo "INFO: Role name was not defined '$ROLE'" && return 1
    (! $(isNaturalNumber $PERM)) && echoInfo "INFO: Invalid permission id '$PERM' " && return 1
    (! $(isNaturalNumber $TIMEOUT)) && TIMEOUT=180
    if ($(isRolePermBlacklisted "$ROLE" "$PERM")) ; then
        echoWarn "WARNING: Role '$ROLE' blacklist already has assigned permission '$PERM'"
    else
        if ($(isRolePermWhitelisted "$ROLE" "$PERM")) ; then
            echoWarn "WARNING: Role '$ROLE' has whitelisted permission '$PERM', attempting to clear..."
            roleClearPermission $ACCOUNT $ROLE $PERM $TIMEOUT
        fi

        sekaid tx customgov role blacklist-permission "$ROLE" "$PERM" --from "$ACCOUNT" --keyring-backend=test --chain-id=$NETWORK_NAME --home=$SEKAID_HOME --fees=100ukex --yes --broadcast-mode=async --log_format=json --output=json | txAwait $TIMEOUT
    fi
}

# createRoleProposal <account> <role-name> <role-description>
# e.g. createRoleProposal validator validator "Role enabling to claim validator seat and perform essential gov. functions"
function createRoleProposal() {
    local ACCOUNT=$1
    local NAME=$2
    local DESCRIPTION=$3
    local TIMEOUT=$4
    ($(isNullOrEmpty $ACCOUNT)) && echoInfo "INFO: Account name was not defined '$1'" && return 1
    ($(isNullOrEmpty $NAME)) && echoInfo "INFO: Invalid role name '$NAME' " && return 1
    (! $(isNaturalNumber $TIMEOUT)) && TIMEOUT=180
    sekaid tx customgov proposal create-role "$NAME" "$DESCRIPTION" --title="Upsert Governance Role '$NAME'" --description="Role description: '$DESCRIPTION'" --from "$ACCOUNT" --keyring-backend=test --chain-id=$NETWORK_NAME --home=$SEKAID_HOME --fees=100ukex --yes --broadcast-mode=async --log_format=json --output=json | txAwait $TIMEOUT
}

# isRoleAssigned <role-name> <address/account>
# eg.: isRoleAssigned sudo test
function isRoleAssigned() {
    local ROLE=$(showRole "$1" 2> /dev/null | jq ".id" 2> /dev/null || echo -n "")
    local ADDRESS=$(showAddress $2)
    
    if (! $(isNaturalNumber $ROLE)) || ($(isNullOrEmpty $ADDRESS)) ; then
        echo "false"
    else
        INDEX=$(showRoles "$ADDRESS" 2> /dev/null | jq ".roleIds | index($ROLE)" 2> /dev/null || echo -n "")
        (! $(isNaturalNumber $INDEX)) && INDEX=$(showRoles "$ADDRESS" 2> /dev/null | jq ".roleIds | index(\"$ROLE\")" 2> /dev/null || echo -n "")
        ($(isNaturalNumber $INDEX)) && echo "true" || echo "false"
    fi
}

# assignRole <account> <role-name> <address> <timeout>
# e.g.: assignRole validator sudo test 180
function assignRole() {
    local ACCOUNT="$1"
    local ROLE=$(showRole "$2" 2> /dev/null | jq ".id" 2> /dev/null || echo -n "")
    local ADDRESS=$(showAddress $3)
    
    local TIMEOUT=$4

    ($(isNullOrEmpty $ACCOUNT)) && echoErr "ERROR: Account name was not defined " && return 1
    ($(isNullOrEmpty $ADDRESS)) && echoErr "ERROR: Assignement address was not defined " && return 1
    (! $(isNaturalNumber $ROLE)) && echoErr "ERROR: Unknown or undefined role '$2'" && return 1
    (! $(isNaturalNumber $TIMEOUT)) && TIMEOUT=180

    if ($(isRoleAssigned "$ROLE" "$ADDRESS")) ; then
        echoWarn "WARNING: Role '$2' was already assigned to account '$ADDRESS'"
    else
        echoInfo "INFO: Adding role '$2' to account '$ADDRESS'"
        sekaid tx customgov role assign "$ROLE" --addr=$ADDRESS --from=$ACCOUNT --keyring-backend=test --chain-id=$NETWORK_NAME --fees=100ukex --yes --log_format=json --broadcast-mode=async --output=json --home=$SEKAID_HOME | txAwait $TIMEOUT
    fi
}

# removeRole <account> <role-name> <address> <timeout>
# e.g.: removeRole validator sudo test 180
function removeRole() {
    local ACCOUNT="$1"
    local ROLE=$(showRole "$2" 2> /dev/null | jq ".id" 2> /dev/null || echo -n "")
    local ADDRESS=$(showAddress $3)
    local TIMEOUT=$4

    ($(isNullOrEmpty $ACCOUNT)) && echoErr "ERROR: Account name was not defined " && return 1
    ($(isNullOrEmpty $ADDRESS)) && echoErr "ERROR: Assignement address was not defined " && return 1
    (! $(isNaturalNumber $ROLE)) && echoErr "ERROR: Unknown or undefined role '$2'" && return 1
    (! $(isNaturalNumber $TIMEOUT)) && TIMEOUT=180

    if (! $(isRoleAssigned "$ROLE" "$ADDRESS")) ; then
        echoWarn "WARNING: Role '$2' was already removed or never assigned to account '$ADDRESS'"
    else
        echoInfo "INFO: Removing role '$2' from account '$ADDRESS'"
        sekaid tx customgov role remove "$ROLE" --addr=$ADDRESS --from=$ACCOUNT --keyring-backend=test --chain-id=$NETWORK_NAME --fees=100ukex --yes --log_format=json --broadcast-mode=async --output=json --home=$SEKAID_HOME | txAwait $TIMEOUT
    fi
}

# whitelistValidators <account> <file-name>
# e.g.: whitelistValidators validator ./whitelist
function whitelistValidators() {
    local ACCOUNT=$1
    local WHITELIST=$2
    local TIMEOUT=$3

    (! $(isNaturalNumber $TIMEOUT)) && TIMEOUT=180
    if [ -f "$WHITELIST" ] ; then 
        echoInfo "INFO: List of validators was found ($WHITELIST)"
        while read key ; do
            key=$(echo "$key" | xargs || echo -n "")
            if ($(isNullOrEmpty "$key")) ; then
                echoWarn "INFO: Invalid key $key"
                continue
            fi

            echoInfo "INFO: Fueling address $WHITELIST with funds from $ACCOUNT"
            sekaid tx bank send $ACCOUNT $key "954321ukex" --keyring-backend=test --chain-id=$NETWORK_NAME --fees 100ukex --yes --log_format=json --broadcast-mode=async --output=json --home=$SEKAID_HOME | txAwait $TIMEOUT

            echoInfo "INFO: Whitelisting '$key' using account '$ACCOUNT'"
            assignRole "$ACCOUNT" validator "$key" "$TIMEOUT" || echoErr "ERROR: Failed to whitelist $key within ${TIMEOUT}s"
        done < $WHITELIST
    elif ($(isKiraAddress $WHITELIST)) ; then
        echoInfo "INFO: Fueling address $WHITELIST with funds from $ACCOUNT"
        sekaid tx bank send $ACCOUNT $WHITELIST "954321ukex" --keyring-backend=test --chain-id=$NETWORK_NAME --fees 100ukex --yes --log_format=json --broadcast-mode=async --output=json --home=$SEKAID_HOME | txAwait $TIMEOUT

        assignRole "$ACCOUNT" validator "$WHITELIST" "$TIMEOUT" || echoErr "ERROR: Failed to whitelist $key within ${TIMEOUT}s"
    else
        echoErr "ERROR: List of validators was NOT found ($WHITELIST)"
    fi
}

# unjailValidators <account> <file-name>
# e.g.: unjailValidators validator ./jailed
# curl -s https://testnet-rpc.kira.network/api/valopers?all=true | jq '.validators | .[] | select(.status=="JAILED")'  | grep -o '".*"' | sed 's/"//g' | grep -o '\bkira1\w*' | uniq > ./jailed
function unjailValidators() {
    local ACCOUNT=$1
    local ADDRESSES=$2
    local TIMEOUT=$3

    (! $(isNaturalNumber $TIMEOUT)) && TIMEOUT=180
    if [ -f "$ADDRESSES" ] ; then 
        echoInfo "INFO: List of validators was found ($ADDRESSES)"
        while read key ; do
            key=$(echo "$key" | xargs || echo -n "")
            if ($(isNullOrEmpty "$key")) ; then
                echoWarn "INFO: Invalid key $key"
                continue
            fi
            echoInfo "INFO: Unjailing '$key' using account '$ACCOUNT'"
            unjail $ACCOUNT "$key" || echoErr "ERROR: Failed to unjail $key within ${TIMEOUT}s"

            echoInfo "INFO: Searching for the last proposal submitted on-chain and voting YES"
            voteYes $(lastProposal) $ACCOUNT  || echoErr "ERROR: Failed to vote yes on the last proposal"
        done < $ADDRESSES
    elif ($(isKiraAddress $ADDRESSES)) ; then
        echoInfo "INFO: Unjailing '$key' using account '$ACCOUNT'"
        unjail $ACCOUNT "$key" || echoErr "ERROR: Failed to unjail $key within ${TIMEOUT}s"

        echoInfo "INFO: Searching for the last proposal submitted on-chain and voting YES"
        voteYes $(lastProposal) $ACCOUNT  || echoErr "ERROR: Failed to vote yes on the last proposal"
    else
        echoErr "ERROR: List of validators was NOT found ($ADDRESSES)"
    fi
}

# clearPermissions <account> <permission> <addresses-file>
# e.g.: clearPermissions validator 29 ./whitelist
function clearPermissions() {
    local ACCOUNT=$1
    local PERMISSION=$2
    local ADDRESSES=$3
    local TIMEOUT=$3

    (! $(isNaturalNumber $TIMEOUT)) && TIMEOUT=180
    if [ -f "$ADDRESSES" ] ; then 
        echoInfo "INFO: List of validators was found ($ADDRESSES)"
        while read key ; do
            key=$(echo "$key" | xargs || echo -n "")
            if ($(isNullOrEmpty "$key")) ; then
                echoWarn "INFO: Invalid key $key"
                continue
            fi
            echoInfo "INFO: Clearing permission '$PERMISSION' for account '$key' using account '$ACCOUNT'"
            clearPermission $ACCOUNT $PERMISSION "$key" $TIMEOUT || echoErr "ERROR: Failed cearing permission '$PERMISSION' for account '$key' using account '$ACCOUNT'"
        done < $ADDRESSES
    elif ($(isKiraAddress $ADDRESSES)) ; then
        echoInfo "INFO: Clearing permission '$PERMISSION' for account '$key' using account '$ACCOUNT'"
        clearPermission $ACCOUNT $PERMISSION "$key" $TIMEOUT || echoErr "ERROR: Failed cearing permission '$PERMISSION' for account '$key' using account '$ACCOUNT'"
    else
        echoErr "ERROR: List of validators was NOT found ($ADDRESSES)"
    fi
}
#showAddress tester3
function isAccount() {
    local ACCOUNT=$(toLower "$1" | xargs || echo -n "")
    local KEYS=$(showKeys 2> /dev/null | jsonParse "" 2> /dev/null || echo "[]")
    readarray -t KEYS_ARR < <(jq -c '.[]' <<< $KEYS)
    for key in "${KEYS_ARR[@]}"; do
      local key_name=$(jq '.name' <<< "$key" | xargs)
      [ "$(toLower "$key_name")" == $ACCOUNT ] && echo "true" && return 0
    done
    echo "false"
}   

function addAccount() {
    local ACCOUNT=$(toLower "$1" | xargs || echo -n "")
    ($(isAccount "$ACCOUNT")) && echoErr "ERROR: Account '$ACCOUNT' already exists" && return 1
    (! $(isAlphanumeric "$ACCOUNT")) && echoErr "ERROR: Account name must be alphanumeric, but got '$ACCOUNT'" && return 1
    sekaid keys add "$ACCOUNT" --keyring-backend=test --home=$SEKAID_HOME --output=json | jq
}

function recoverAccount() {
    local ACCOUNT=$(toLower "$1" | xargs || echo -n "")
    local MNEMONIC=$(toLower "$2" | xargs || echo -n "")
    ($(isAccount "$ACCOUNT")) && echoErr "ERROR: Account '$ACCOUNT' already exists" && return 1
    (! $(isAlphanumeric "$ACCOUNT")) && echoErr "ERROR: Account name must be alphanumeric, but got '$ACCOUNT'" && return 1
    yes "$MNEMONIC" | sekaid keys add $ACCOUNT --keyring-backend=test --home=$SEKAID_HOME --recover --output=json | jq
}

function deleteAccount() {
    local ACCOUNT=$(toLower "$1" | xargs || echo -n "")
    if (! $(isAccount "$ACCOUNT")) ; then 
        echoWarn "WARNING: Account '$ACCOUNT' was already deleted or never existed"
    else
        sekaid keys delete "$ACCOUNT" --force --yes --keyring-backend=test --home=$SEKAID_HOME --output=json | jq
    fi
}

# allow to execute finctions directly from file
if declare -f "$1" > /dev/null ; then
  # call arguments verbatim
  "$@"
fi

# enableCustody
function enableCustody() {
  local FROM=$1
  local MODE=$2
  local PASSWORD=$3
  local LIMITS=$4
  local WHITELIST=$5
  local FEE_AMOUNT=$6
  local FEE_DENOM=$7
  local OKEY=$8
  local NKEY=$9

  sekaid tx custody create $MODE $PASSWORD $LIMITS $WHITELIST --from=$FROM --keyring-backend=test --chain-id=$NETWORK_NAME --fees="${FEE_AMOUNT}${FEE_DENOM}" --output=json --yes --home=$SEKAID_HOME --okey=$OKEY --nkey=$NKEY | txAwait 180
}
# enableCustody
function disableCustody() {
  local FROM=$1
  local FEE_AMOUNT=$2
  local FEE_DENOM=$3
  local OKEY=$4

  sekaid tx custody disable --from=$FROM --keyring-backend=test --chain-id=$NETWORK_NAME --fees="${FEE_AMOUNT}${FEE_DENOM}" --output=json --yes --home=$SEKAID_HOME --okey=$OKEY | txAwait 180
}

# addCustodians
function addCustodians() {
  local FROM=$1
  local ADDRESSES=$2
  local FEE_AMOUNT=$3
  local FEE_DENOM=$4
  local OKEY=$5
  local NKEY=$6

  sekaid tx custody custodians add "$ADDRESSES" --from=$FROM --keyring-backend=test --chain-id=$NETWORK_NAME --fees="${FEE_AMOUNT}${FEE_DENOM}" --output=json --yes --home=$SEKAID_HOME --okey=$OKEY --nkey=$NKEY | txAwait 180
}

# dropCustodians
function dropCustodians() {
  local FROM=$1
  local FEE_AMOUNT=$2
  local FEE_DENOM=$3
  local OKEY=$4
  local NKEY=$5

  sekaid tx custody custodians drop --from=$FROM --keyring-backend=test --chain-id=$NETWORK_NAME --fees="${FEE_AMOUNT}${FEE_DENOM}" --output=json --yes --home=$SEKAID_HOME --okey=$OKEY --nkey=$NKEY | txAwait 180
}

# e.g. getCustodyInfo
function getCustodyKey() {
  local FROM=$(showAddress $1)
  local RESULT=$(sekaid query custody get $FROM --output=json --home=$SEKAID_HOME 2> /dev/null | jsonParse "custody_settings.key" 2> /dev/null || echo -n "")

  echo $RESULT
}

# e.g. getCustodyInfo
function getCustodyInfo() {
  local FROM=$(showAddress $1)
  local RESULT=$(sekaid query custody get $FROM --output=json --home=$SEKAID_HOME 2> /dev/null | jsonParse "" 2> /dev/null || echo -n "")

  echo $RESULT
}

# e.g. getCustodyWhitelist
function getCustodians() {
  local FROM=$(showAddress $1)
  local RESULT=$(sekaid query custody custodians get $FROM --output=json --home=$SEKAID_HOME 2> /dev/null | jsonParse "custody_custodians" 2> /dev/null || echo -n "")

  echo $RESULT
}

# e.g. getCustodyWhitelist
function getCustodyPool() {
  local FROM=$(showAddress $1)
  local RESULT=$(sekaid query custody custodians pool $FROM --output=json --home=$SEKAID_HOME 2> /dev/null | jsonParse "transactions.record" 2> /dev/null || echo -n "")

  echo $RESULT
}

function getCustodyPoolVotes() {
  local FROM=$(showAddress $1)
  local HASH=$2

  local RESULT=$(sekaid query custody custodians pool $FROM --output=json --home=$SEKAID_HOME 2> /dev/null | jsonParse "transactions.record.$HASH.votes" 2> /dev/null || echo -n "")

  echo $RESULT
}

# e.g. sendTokens faucet kiraXXX...XXX 1000 ukex 100 ukex
function custodySendTokens() {
    local SOURCE=$1
    local DESTINATION=$(showAddress $2)
    local AMOUNT="$3"
    local DENOM="$4"
    local FEE_AMOUNT="$5"
    local FEE_DENOM="$6"
    local REWARD_AMOUNT="$7"
    local REWARD_DENOM="$8"
    local PASSWORD="$9"
    local RESULT=""

    ($(isNullOrEmpty $FEE_AMOUNT)) && FEE_AMOUNT=100
    ($(isNullOrEmpty $FEE_DENOM)) && FEE_DENOM="ukex"

    RESULT=$(sekaid tx custody send $SOURCE $DESTINATION "${AMOUNT}${DENOM}" $PASSWORD --reward="${REWARD_AMOUNT}${REWARD_DENOM}" --keyring-backend=test --chain-id=$NETWORK_NAME --fees "${FEE_AMOUNT}${FEE_DENOM}" --output=json --yes --home=$SEKAID_HOME 2> /dev/null | txAwait2 180 2> /dev/null || echo -n "" )

    echo "${RESULT,,}"
}

function approveTransaction() {
    local FROM=$1
    local ADDRESS=$(showAddress $2)
    local HASH=$3
    local FEE_AMOUNT="$4"
    local FEE_DENOM="$5"

    sekaid tx custody approve --from=$FROM $ADDRESS $HASH --keyring-backend=test --chain-id=$NETWORK_NAME --fees "${FEE_AMOUNT}${FEE_DENOM}" --output=json --yes --home=$SEKAID_HOME | txAwait 180
}

function declineTransaction() {
    local FROM=$1
    local ADDRESS=$(showAddress $2)
    local HASH=$3
    local FEE_AMOUNT="$4"
    local FEE_DENOM="$5"

    sekaid tx custody decline --from=$FROM $ADDRESS $HASH --keyring-backend=test --chain-id=$NETWORK_NAME --fees "${FEE_AMOUNT}${FEE_DENOM}" --output=json --yes --home=$SEKAID_HOME | txAwait 180
}

function passwordConfirmTransaction() {
    local FROM=$1
    local ADDRESS=$(showAddress $2)
    local HASH=$3
    local PASSWORD=$4
    local FEE_AMOUNT="$5"
    local FEE_DENOM="$6"

    sekaid tx custody confirm --from=$FROM $ADDRESS $HASH $PASSWORD --keyring-backend=test --chain-id=$NETWORK_NAME --fees "${FEE_AMOUNT}${FEE_DENOM}" --output=json --yes --home=$SEKAID_HOME | txAwait 180
}