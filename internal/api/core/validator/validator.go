package validator

import (
	"errors"
	"regexp"
)

func ValidateEmail(email string) error {
    if email == "" {
        return errors.New("Email cannot be empty.")
    }

    valid, err := regexp.Match(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, []byte(email))
    if err != nil {
        return err
    }

    if !valid {
        return errors.New("Invalid email format.")
    }

    return nil
}

func ValidateUsername(username string) error {
    if username == "" {
        return errors.New("Username cannot be empty.")
    }

    if len(username) > 40 {
        return errors.New("Username maximum length is 16 characters.")
    }

    startWithLetter, err := regexp.Match(`^([A-Za-z_])`, []byte(username))
    if err != nil {
        return err
    }
    if !startWithLetter {
        return errors.New("Username must start with letter or underscore.")
    }

    onlyAllowCharacter, err := regexp.Match(`[A-Za-z0-9_\/]`, []byte(username))
    if err != nil {
        return err
    }
    if !onlyAllowCharacter {
        return errors.New("Username can contain only letter, number and underscore.")
    }

    return nil
}

func ValidateDisplayname(displayname string) error {
    valid, err := regexp.Match(`^.{1,70}$`, []byte(displayname))
    if err != nil {
        return err
    }

    if !valid {
        return errors.New("Invalid display name format.")
    }

    return nil
}

func ValidatePassword(password string) error {
    if password == "" {
        return errors.New("Password cannot be empty.")
    }

    {
        valid, err := regexp.Match(`^.{8,255}`, []byte(password))
        if err != nil {
            return err
        }
        if !valid {
            return errors.New("Invalid password format.")
        }
    }

    {
        valid, err := regexp.Match(`^(.*[A-Za-z])`, []byte(password))
        if err != nil {
            return err
        }
        if !valid {
            return errors.New("Invalid password format.")
        }
    }

    {
        valid, err := regexp.Match(`^(.*[\d])`, []byte(password))
        if err != nil {
            return err
        }
        if !valid {
            return errors.New("Invalid password format.")
        }
    }


    return nil
}

