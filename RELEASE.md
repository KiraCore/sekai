Features:

- [TYPO] arraw => array
- implement token inflation distribution
- implementation of reward cap per staking token
- Add staking token and stake min validation on delegation, validation on total reward cap
- Update CLI commands for token rates msg and proposal
- CLI test for reward cap exceeding total limitation
- Resolve issues on autocompounding
- implement delegators cap
- implement pushout when 10x of minimum
- add staking pool delegators query
- automatic unregister of delegator on a pool if not enough stake available
- prevent delegation to not an active validator
- add tests on changes
- fixed github workflows, now all release branches have diffrent names then version tags
- speedup of integraiton tests
