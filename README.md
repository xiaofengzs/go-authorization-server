http://localhost:8080/login

http://localhost:8080/oauth2/auth?client_id=my-client&response_type=code&state=123455611111

ory_ac_koogQyahUg6o95IuAsID8Bta5mcfZG99NB1ieT61JbU.E-A67YFhLJ1p-Us34n1Vne2-lUWeF8XMPzepkgiIAMM

ory_ac_ErjJ9fbR5UOks9FT2Pa2t4MdBYkVxdHr9gIz_28up58.xCkMNuSwZmhnQ4hOXXOh7Cl9YZ5B-srUJpvjrlUy0TE

curl -X POST \
  http://localhost:8080/oauth2/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "code=ory_ac_ErjJ9fbR5UOks9FT2Pa2t4MdBYkVxdHr9gIz_28up58.xCkMNuSwZmhnQ4hOXXOh7Cl9YZ5B-srUJpvjrlUy0TE" \
  -d "client_id=my-client" \
  -d "client_secret=foobar" \
  -d "redirect_uri=localhost:8080/callback" \
  -d "grant_type=authorization_code"
