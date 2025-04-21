Overview
========

A very short demonstration of a SAML service provider.

To test:

1) Run keycloak with: `docker run -p 8080:8080 -e KEYCLOAK_ADMIN=admin -e KEYCLOAK_ADMIN_PASSWORD=admin quay.io/keycloak/keycloak:latest start-dev`

2) Browse to http://localhost:8080 and login with `admin`, `admin`.

3) Create a keycloak realm for testing. Name it `example-oidc`.

4) In keycloak create a `Client`. Inside `demo-realm`, go to `Clients`, click `Create Client`, use the following values:

```
Client ID: demo-client
Client Type: OpenID Connect
Root URL: http://localhost:3000
```

Click `Save`.

Under Settings:

Enable `Standard Flow` and `Client Authentication`.

Add Redirect URI: http://localhost:3000/callback

Optional: Set Web Origins: http://localhost:3000

Go to `Credentials` tab and copy the Client Secret.


5) In keycloak, add a user: Users > Add user. Username: `example`. Create. Give the user a password from the `Credentials` tab.

6) Run `main.go` with the following environment variables:
```bash
export CLIENT_ID=demo-client
export CLIENT_SECRET=your-client-secret-here
export REDIRECT_URL=http://localhost:3000/callback
export KEYCLOAK_URL=http://localhost:8080/realms/demo-realm
go run .
```

7) After the server starts, navigate to `http://localhost:3000` in your browser and click the link to use Keycloak to login as your test user.
