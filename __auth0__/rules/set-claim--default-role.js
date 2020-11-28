function SetDefaultRoles(user, context, callback) {
    const namespace = 'https://client.devpie.io/claims/roles';

    context.idToken[namespace] = [];

    if (user.email.indexOf('@devpie.io') !== -1) {
        context.idToken[namespace].push('employee');
    } else {
        context.idToken[namespace].push('user');
    }

    callback(null, user, context);
}