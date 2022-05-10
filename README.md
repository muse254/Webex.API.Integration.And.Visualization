# Integration Visualization with Webex API

The program code is written in Golang and HTML. The frontend is rendered as a static site.

The application is focused on fetching `Meeting Analytics Quality` data, this data is only accessible for `1 request every 5 minutes`. Therefore, there is need to have persitance in our application. The program uses an SQlite database for this just to allow for the fault tolerance on the part of the client making multiple requests. SQlite was chosen because it is a simple, easy to use and requires little to no setup to have it working. The data from the metrics is not interpreted or manipulated in any way, only a `read` and `write` mechanism exists.

## Authorization

The program requires that the client creates an integration with Webex. During the integration creation the scopes needed will be: `meeting:schedules_read` and `analytics:read_all`.
- `meeting:schedules_read` is needed to fetch all meetings that have been conducted for the account.
- `analytics:read_all` is needed to fetch all meeting room activity, this includes quality metrics.

On success the client will be provided with the `client_id` and `client_secret` that will be used to initiate the OAuth Flow. If the `scope` of the application changes, the client secret will also change and the previous `client_secret` is rendered invalid(important to note). In this event the Webex Integration page for your account will have a new `client_secret`.

The OAuth Flow is used to request the permissions asked for through tradition login for the client. If successful it triggers the server with a redirect to the client's auth redirect URI. The redirect URI points to our server and it also needs to be publicly accessible and able to accept the request. This is possible by hosting our codebase at a Cloud provider like `Amazon Web Services`.

The redirect request from the OAuth Flow after succesful login comes with a `code` parameter that is used to fetch an `access token` with a `refresh token`. To request for this token the client also needs to provide the `code`; and also the `client_id` and `client_secret` that were provided during the integration creation. With this bit out of the way our final response is a JSON object with the following fields:
```json
{
    "access_token": "...",
    "expires_in": 0000,
    "refresh_token": "...",
    "refresh_token_expires_in": 0000,
}
```

The `access_token` is used to request for any resources under the scope it was approved for in the Webex API. The `refresh_token` is used to request a new `access_token` after it has expired.

## APIs

Our integration is focused on checking on the meeting analytics quality and is minimal in the number of APIs it uses.
1. First, we use [List Meetings API](https://developer.webex.com/docs/api/v1/meetings/list-meetings) to request for all meetings available to us. With this we can use their `meeting_id`s in the next API.
2. Lastly, we use []