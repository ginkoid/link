# link

Super simple redirect service! Powers [gnk.io](https://gnk.io).

Downloads redirects from a gist.

## deploying your own

Clone the repo to your machine.

Make a gist with a single file called `link`, with JSON content that looks like this:

```json
{
  "/path/for/redirect": { "to": "https://example.com" },
  "/":                  { "to": "https://example.org" }
}
```

[create a GitHub OAuth app](https://github.com/settings/applications/new). This won't be used to authenticate any users, just to have higher rate limits for the GitHub API.

Edit the `now.json` file to change the `alias`es to your domains. Or, remove them to have a default [now.sh](https://zeit.co/now) url.

Use the [now CLI](https://www.npmjs.com/package/now) to create these secrets:

* `app-link-github-gist-id`: the ID gist you created
* `app-link-github-client-id`: the client ID of the GitHub OAuth app you created
* `app-link-github-client-secret`: the client secret of the GitHub OAuth app you created

Run `now` to deploy the app.
