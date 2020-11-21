function SetUserId(user, context, callback) {
    const namespace = 'https://client.devpie.io/claims/';
    user.app_metadata = user.app_metadata || {};

    if(user.app_metadata.id) {
        context.accessToken[namespace + 'user_id'] = user.app_metadata.id;
    }

    callback(null, user, context);
}