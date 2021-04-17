export interface ReqContext {
  userId: string;
  traceId: string;
}

export interface Auth0Claims {
  "https://client.devpie.io/claims/is_new": boolean;
  "https://client.devpie.io/claims/roles": "freelancer" | "client";
  sub: `auth0|${string}`;
}

export interface Auth0User {
  auth0Id: string;
  email: string;
  emailVerified: boolean;
  firstName: string;
  lastName: string;
  picture: string;
  locale: string;
}

export interface DBUser {
  user_id: string;
  auth0_id: string;
  email: string;
  email_verified: boolean;
  first_name: string;
  last_name: string;
  picture: string;
  locale: string;
}

export interface User {
  userId: string;
  auth0Id: string;
  email: string;
  emailVerified: boolean;
  firstName: string;
  lastName: string;
  picture: string;
  locale: string;
}
