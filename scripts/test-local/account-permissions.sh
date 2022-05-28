#!/usr/bin/env bash
# To run test locally: make network-start && ./scripts/test-local/account-permissions.sh
set -e
set -x
. /etc/profile
. ./scripts/sekai-env.sh

TEST_NAME="ACCOUNT-PERMISSIONS" && timerStart $TEST_NAME
echoInfo "INFO: $TEST_NAME - Integration Test - START"

echoInfo "INFO: Whitelisting & veryfying individual account permissions..."

# NOTE: We can't remove 'PermSetPermissions' otherwise sudo power would be lost

declare -a perms=(
    "PermZero"
    "PermClaimValidator"
    "PermClaimCouncilor"
    "PermWhitelistAccountPermissionProposal"
    "PermVoteWhitelistAccountPermissionProposal"
    "PermUpsertTokenAlias"
    "PermChangeTxFee"
    "PermUpsertTokenRate"
    "PermUpsertRole"
    "PermCreateUpsertDataRegistryProposal"
    "PermVoteUpsertDataRegistryProposal"
    "PermCreateSetNetworkPropertyProposal"
    "PermVoteSetNetworkPropertyProposal"
    "PermCreateUpsertTokenAliasProposal"
    "PermVoteUpsertTokenAliasProposal"
    "PermCreateSetPoorNetworkMessagesProposal"
    "PermVoteSetPoorNetworkMessagesProposal"
    "PermCreateUpsertTokenRateProposal"
    "PermVoteUpsertTokenRateProposal"
    "PermCreateUnjailValidatorProposal"
    "PermVoteUnjailValidatorProposal"
    "PermCreateRoleProposal"
    "PermVoteCreateRoleProposal"
    "PermCreateTokensWhiteBlackChangeProposal"
    "PermVoteTokensWhiteBlackChangeProposal"
    "PermCreateResetWholeValidatorRankProposal"
    "PermVoteResetWholeValidatorRankProposal"
    "PermCreateSoftwareUpgradeProposal"
    "PermVoteSoftwareUpgradeProposal"
    "PermSetClaimValidatorPermission"
    "PermCreateSetProposalDurationProposal"
    "PermVoteSetProposalDurationProposal"
    "PermBlacklistAccountPermissionProposal"
    "PermVoteBlacklistAccountPermissionProposal"
    "PermRemoveWhitelistedAccountPermissionProposal"
    "PermVoteRemoveWhitelistedAccountPermissionProposal"
    "PermRemoveBlacklistedAccountPermissionProposal"
    "PermVoteRemoveBlacklistedAccountPermissionProposal"
    "PermWhitelistRolePermissionProposal"
    "PermVoteWhitelistRolePermissionProposal"
    "PermBlacklistRolePermissionProposal"
    "PermVoteBlacklistRolePermissionProposal"
    "PermRemoveWhitelistedRolePermissionProposal"
    "PermVoteRemoveWhitelistedRolePermissionProposal"
    "PermRemoveBlacklistedRolePermissionProposal"
    "PermVoteRemoveBlacklistedRolePermissionProposal"
    "PermAssignRoleToAccountProposal"
    "PermVoteAssignRoleToAccountProposal"
    "PermUnassignRoleFromAccountProposal"
    "PermVoteUnassignRoleFromAccountProposal"
    "PermRemoveRoleProposal"
    "PermVoteRemoveRoleProposal")

for p in "${perms[@]}" ; do
    echoInfo "INFO: Whitelisting or blacklisting & checking account permission '$p'..."

    PERM_ID="${!p}"
    PERM_CHECK=$(isPermWhitelisted validator $PERM_ID)
    [ "${PERM_CHECK,,}" != "false" ] && echoErr "ERROR: Expected account validator to NOT have permission '$p' ($PERM_ID) before its whitelisted" && exit 1
    PERM_CHECK=$(isPermBlacklisted validator $PERM_ID)
    [ "${PERM_CHECK,,}" != "false" ] && echoErr "ERROR: Expected account validator to NOT have blacklisted permission '$p' ($PERM_ID) before its blacklisted" && exit 1

    # here we are testing whitelisting & blacklisting, every even permission is whitelisted, every odd is blacklisted
    if [[ $(( $PERM_ID % 2 )) == 0 ]]; then
        whitelistPermission validator $PERM_ID validator

        PERM_CHECK=$(isPermWhitelisted validator $PERM_ID)
        [ "${PERM_CHECK,,}" != "true" ] && echoErr "ERROR: Expected account validator to have permission '$p' ($PERM_ID) after it was whitelisted" && exit 1
        echoInfo "INFO: Success, whitelised permission '$p'..."
    else
        blacklistPermission validator $PERM_ID validator

        PERM_CHECK=$(isPermBlacklisted validator $PERM_ID)
        [ "${PERM_CHECK,,}" != "true" ] && echoErr "ERROR: Expected account validator to have blacklisted permission '$p' ($PERM_ID) after it was blacklisted" && exit 1
        echoInfo "INFO: Success, blacklisted permission '$p'..."
    fi

    clearPermission validator $PERM_ID validator
    PERM_CHECK=$(isPermWhitelisted validator $PERM_ID)
    [ "${PERM_CHECK,,}" != "false" ] && echoErr "ERROR: Expected account validator to NOT have permission '$p' ($PERM_ID) after it was cleared" && exit 1
    PERM_CHECK=$(isPermBlacklisted validator $PERM_ID)
    [ "${PERM_CHECK,,}" != "false" ] && echoErr "ERROR: Expected account validator to NOT have blacklisted permission '$p' ($PERM_ID) after it was cleared" && exit 1
done

echoInfo "INFO: $TEST_NAME - Integration Test - END, elapsed: $(prettyTime $(timerSpan $TEST_NAME))"
