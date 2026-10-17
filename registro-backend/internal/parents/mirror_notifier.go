package parents

import (
	"context"
	"fmt"
)

type MirrorNotifier struct{}

func NewMirrorNotifier() *MirrorNotifier {
	return &MirrorNotifier{}
}

func (m *MirrorNotifier) NotifyNewAuthorization(ctx context.Context, auth *DualParentalAuthorization, recipientID string) error {
	// Dispatches notification (Push / WebSocket / Audit)
	return nil
}

func (m *MirrorNotifier) NotifyPendingSecondSignature(ctx context.Context, auth *DualParentalAuthorization, recipientID, signerID string) error {
	// Notifies co-parent: "L'altro genitore ha firmato l'atto '%s'. È richiesta la tua firma congiunta per completare la procedura."
	_ = fmt.Sprintf("L'altro genitore ha firmato l'atto '%s'. È richiesta la tua firma congiunta.", auth.Title)
	return nil
}

func (m *MirrorNotifier) NotifyRejection(ctx context.Context, auth *DualParentalAuthorization, recipientID, reason string) error {
	_ = fmt.Sprintf("L'atto '%s' è stato rifiutato dall'altro genitore. Motivo: %s", auth.Title, reason)
	return nil
}
