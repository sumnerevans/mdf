# mdf (Mutt Display Filter)

A custom mutt display filter. It has two components: a daemon for minifying URLs
and an executable to use as the display filter.

## Features

* URL shortening (works by running a redirecter daemon which serves redirect
  pages)
* Email formatting normalization
* Date header normalization to local timezone
* Git diff highlighting

## Redirecter Usage

Note that it is not guaranteed to generate unique IDs, but that's fine since
they are random enough.

### `POST /new`

Create a new redirect page. Send the URL in the body of the request.

Returns an ID to use below.

### `GET /{id}`

Serve a redirect page.
