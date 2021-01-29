function SetAssignedRoles(user, context, callback) {
    const namespace = 'https://client.devpie.io/claims/roles';

    const assignedRoles = (context.authorization || {}).roles;

    context.idToken[namespace] = assignedRoles;
    context.accessToken[namespace] = assignedRoles;

    callback(null, user, context);
}