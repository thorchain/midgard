package binance

import "time"

// TxResponse
type TxResponse struct {
	Hash   string `json:"hash"`
	Height string `json:"height"`
	Tx     string `json:"tx"`
}

// FullNodeTxResp full node response
type FullNodeTxResp struct {
	Result TxResponse `json:"result"`
}

// BlockResponse represent the block
// since we only need the timestamp, thus we define it ourselves
type BlockResponse struct {
	Result struct {
		Block struct {
			Header struct {
				Time time.Time `json:"time"`
			} `json:"header"`
		} `json:"block"`
	} `json:"result"`
}
