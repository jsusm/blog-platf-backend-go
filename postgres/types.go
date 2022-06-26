package postgres

import (
  "time"
  "golang.org/x/crypto/bcrypt"
)

type Post struct {
  ID        int       `ksql:"post_id" json:"id"`
  Title     string    `ksql:"title" json:"title"`
  Desc      string    `ksql:"desc" json:"desc"`
  Content   string    `ksql:"content" json:"content"`
  CreatedAt time.Time `ksql:"created_at" json:"created_at"`
  UpdatedAt time.Time `ksql:"updated_at" json:"updated_at"`
  Slug      string    `ksql:"slug" json:"slug"`
  SmImg     string    `ksql:"sm_img" json:"smImg"`
  LgImg     string    `ksql:"lg_img" json:"lgImg"`
  UserId    string    `ksql:"user_id" json:"userID"`
}

type User struct {
  ID       int    `ksql:"user_id"`
	Name     string `ksql:"name"`
	password string `ksql:"password"`
	Img      string `ksql:"img"`
}

// Encrypt the password and set to the struct
func (u User) EncryptAndSetPassword(newPassword string)(error){
  pass, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
  if err != nil {
    return err
  }
  u.password = string(pass)
  return nil
}

// Check if the given password is correct
// return true on success and false on failure
func (u User) ValidPassword(password string)(bool){
  err := bcrypt.CompareHashAndPassword([]byte(u.password), []byte(password))
  return err == nil
}
