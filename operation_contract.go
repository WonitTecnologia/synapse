package synapse

import (
	"bytes"
	"encoding/json"
)

// OperationContract descreve o significado de uma operação executada por uma API tool.
// O contrato é configurado pelo tenant e não é inferido pelo nome da ferramenta.
type OperationContract struct {
	Domain   string           `json:"domain"`
	Verb     string           `json:"verb"`
	Entity   string           `json:"entity"`
	Evidence EvidenceContract `json:"evidence"`
}

// EvidenceContract indica onde a resposta comprova sucesso ou disponibilidade.
type EvidenceContract struct {
	SuccessPath      string `json:"success_path,omitempty"`
	SuccessValue     any    `json:"success_value"`
	AvailabilityPath string `json:"availability_path,omitempty"`
	AvailableWhen    string `json:"available_when,omitempty"`
	ProtocolPath     string `json:"protocol_path,omitempty"`
}

func (c OperationContract) MarshalJSON() ([]byte, error) {
	if c.Domain == "" && c.Verb == "" && c.Entity == "" && c.Evidence.isZero() {
		return []byte("{}"), nil
	}
	type contractAlias OperationContract
	return json.Marshal(contractAlias(c))
}

func (c *OperationContract) UnmarshalJSON(data []byte) error {
	type contractAlias OperationContract
	var decoded contractAlias
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&decoded); err != nil {
		return err
	}
	*c = OperationContract(decoded)
	return nil
}

func (e EvidenceContract) MarshalJSON() ([]byte, error) {
	payload := map[string]any{}
	if e.SuccessPath != "" {
		payload["success_path"] = e.SuccessPath
	}
	if e.SuccessValue != nil {
		payload["success_value"] = e.SuccessValue
	}
	if e.AvailabilityPath != "" {
		payload["availability_path"] = e.AvailabilityPath
	}
	if e.AvailableWhen != "" {
		payload["available_when"] = e.AvailableWhen
	}
	if e.ProtocolPath != "" {
		payload["protocol_path"] = e.ProtocolPath
	}
	return json.Marshal(payload)
}

func (e EvidenceContract) isZero() bool {
	return e.SuccessPath == "" && e.SuccessValue == nil && e.AvailabilityPath == "" && e.AvailableWhen == "" && e.ProtocolPath == ""
}
