- ! 代表不允許Null


fragment userData on User {
id
name
email
}

query getUser($id: ID!) {
user(id: $id) {
...userData
}
}

query getUsers {
users {
...userData
}
}