package transaction

type EntryID string
type AccountID string
type TransactionID string

type StateMachineStatus map[TransactionStatus][]TransactionStatus

var validStateMachine = StateMachineStatus{
	Pending:   []TransactionStatus{Completed, Failed},
	Completed: []TransactionStatus{},
	Failed:    []TransactionStatus{},
}

func (s StateMachineStatus) CanTransition(from, to TransactionStatus) bool {
	allowed, ok := s[from]
	if !ok {
		return false
	}

	for _, status := range allowed {
		if status == to {
			return true
		}
	}

	return false
}
