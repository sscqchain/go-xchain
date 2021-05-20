package gov

import sdk "gitee.com/xchain/go-xchain/types"

var _ ProposalContent = (*SoftwareUpgradeProposal)(nil)

type SoftwareUpgradeProposal struct {
	Proposal
	ProtocolDefinition sdk.ProtocolDefinition `json:"protocol_definition"`
}

func (sp SoftwareUpgradeProposal) ProposalType() sdk.ProtocolDefinition {
	return sp.ProtocolDefinition
}

func (sp SoftwareUpgradeProposal) GetProtocolDefinition() sdk.ProtocolDefinition {
	return sp.ProtocolDefinition
}
func (sp *SoftwareUpgradeProposal) SetProtocolDefinition(upgrade sdk.ProtocolDefinition) {
	sp.ProtocolDefinition = upgrade
}
