package dto

type TransactionRequest struct {
	SenderPrivateKey          *string `json:"sender_private_key"`
	SenderBlockchainAddress   *string `json:"sender_blockchain_address"`
	ReceiverBlockchainAddress *string `json:"receiver_blockchain_address"`
	SenderPublicKey           *string `json:"sender_public_key"`
	Value                     *string `json:"value"`
}

// add validate field not nill gitu ntar aja

func (t *TransactionRequest) Validate() bool {
	if t.SenderPrivateKey == nil || t.SenderBlockchainAddress == nil || t.ReceiverBlockchainAddress == nil || t.SenderPublicKey == nil || t.Value == nil {
		return false
	}
	return true
}
