/*
Description

Auth0 allows you to define various types of applications, for example, under applications one can define a
single page application, a regular web application or a machine to machine application.

Applications communicate with a user defined API. Auth0 exposes a management API for completing administrative
tasks you would normally do manually in the Auth0 web platform. To connect to applications and APIs you need credentials.
Credentials are found in the Auth0's platform.

Credentials

The following illustrates a list of credentials required by the auth package:

  // Credentials for single page application to API authorization

  API_WEB_AUTH_DOMAIN=
  API_WEB_AUTH_AUDIENCE=

  // Credentials for generating access tokens inside tests

  API_WEB_AUTH_TEST_CLIENT_ID=
  API_WEB_AUTH_TEST_CLIENT_SECRET=

  // Credentials for communicating with the Auth0 management API

  API_WEB_AUTH_M_2_M_CLIENT=
  API_WEB_AUTH_M_2_M_SECRET=
  API_WEB_AUTH_MAPI_AUDIENCE=

The naming convention used is specific to the http://github.com/ardanlabs/conf package.

References

https://auth0.com/blog/authentication-in-golang/

https://auth0.com/docs/flows/call-your-api-using-resource-owner-password-flow

https://auth0.com/blog/using-m2m-authorization/

https://auth0.com/docs/api/management/v2

*/
package auth0
