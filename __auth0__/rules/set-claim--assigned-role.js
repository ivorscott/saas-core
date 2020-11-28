function SetAssignedRoles(user, context, callback) {
    const namespace = 'https://client.devpie.io/claims/roles';

    let idTokenClaims = context.idToken || {};

    const assignedRoles = (context.authorization || {}).roles;
    const defaultRoles = idTokenClaims[namespace] ? idTokenClaims[namespace] : [];
    const combinedRoles = [...assignedRoles, ...defaultRoles];

    context.idToken[namespace] = combinedRoles;
    context.accessToken[namespace] = combinedRoles;

    callback(null, user, context);
}