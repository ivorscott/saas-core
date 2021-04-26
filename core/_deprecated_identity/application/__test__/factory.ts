import camelcaseKeys from "camelcase-keys";
import { User, DBUser, Auth0User } from "../identity";
import { v4 as uuidV4 } from "uuid";

const getFakeUsersRecords = (
  { ...applyPropsToFirstItem }: Partial<DBUser> = {},
  { ...applyPropsToSecondItem }: Partial<DBUser> = {},
  { ...applyPropsToThirdItem }: Partial<DBUser> = {},
): DBUser[] => [
  {
    user_id: uuidV4(),
    auth0_id: "google-oauth2|ABC123",
    email: "example@devpie.io",
    email_verified: true,
    first_name: "Adam",
    last_name: "Smith",
    locale: "en",
    picture: "https://image.com/myimage",
    ...applyPropsToFirstItem,
  },
  {
    user_id: uuidV4(),
    auth0_id: "google-oauth2|XYZ123",
    email: "example@devpie.io",
    email_verified: true,
    first_name: "Marie",
    last_name: "Parker",
    locale: "en",
    picture: "https://image.com/myimage",
    ...applyPropsToSecondItem,
  },
  {
    user_id: uuidV4(),
    auth0_id: "google-oauth2|EFG123",
    email: "example@gmail.com",
    email_verified: true,
    first_name: "Dan",
    last_name: "Wells",
    locale: "en",
    picture: "https://image.com/myimage",
    ...applyPropsToThirdItem,
  },
];

const getFakeAuth0Users = (list: DBUser[]) => {
  const users = (camelcaseKeys(list) as unknown) as User[];
  // skip "userId"
  return users.map(({ userId, ...user }) => {
    return {
      ...user,
    };
  }) as Auth0User[];
};

const formatUserRecords = (records: DBUser[]) =>
  (camelcaseKeys(records) as unknown) as User[];

export { formatUserRecords, getFakeAuth0Users, getFakeUsersRecords };
