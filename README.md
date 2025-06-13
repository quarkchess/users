# auth - user authentication

This service is responsible for handing out JWTs to users who sign in, and also JWTs to
users who _don't_ sign in (e.g. anons). Also, other services can use this service to
cross-check tokens, look up users, etc.
