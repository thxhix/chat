package chat

const (
	maxMembers = 1000
)

func ValidateCreateRequest(cType int8, members []string) error {
	err := ValidateMembers(members)
	if err != nil {
		return err
	}

	err = ValidateDirect(cType, members)
	if err != nil {
		return err
	}
	return nil
}

func ValidateMembers(m []string) error {
	if len(m) <= 0 {
		return ErrEmptyMembersProvided
	}
	if len(m) > maxMembers {
		return ErrTooManyMembersProvided
	}
	return nil
}

func ValidateDirect(cType int8, members []string) error {
	if cType < 1 || cType > 2 {
		return ErrBadDirectTypeProvided
	}
	if cType == 1 && len(members) > 2 {
		return ErrPrivateDirectTooManyMembers
	}
	return nil
}
