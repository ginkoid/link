# link

Super simple redirect service! Downloads redirects from a gist.

## deploying your own

Clone the repo to your machine.

[Make a secret GitHub gist](https://gist.github.com/new) with a single file called `link`, with JSON content that looks like this:

```json
{
  "/path/for/redirect": { "to": "https://example.com" },
  "/":                  { "to": "https://example.org" }
}
```

The paths in the gist should be lowercase. Request paths are lowercased before lookup in the map of redirects.

[Create a GitHub OAuth app](https://github.com/settings/applications/new). This won't be used to authenticate any users, just to have higher rate limits for the GitHub API.

Note that if you have more than ~45 replications of the app, the GitHub ratelimit will be hit, and updates to redirects won't be deployed until the ratelimit resets.

Edit the `now.json` file to change the `alias`es to your domains. Or, remove them to have a default [now.sh](https://zeit.co/now) url.

Use the [now CLI](https://www.npmjs.com/package/now) to create secrets:

```sh
now secret create secret-name secret-value
```

Create these secrets:

* `app-link-github-gist-id`: the ID of the gist you created
* `app-link-github-client-id`: the client ID of the GitHub OAuth app you created
* `app-link-github-client-secret`: the client secret of the GitHub OAuth app you created

Deploy the app.

```sh
now
```
