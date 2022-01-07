const location = window.location;

const dev = {
    app: {
        version: process.env.REACT_APP_VERSION,
    },
    api: {
        wsUrl: 'ws://localhost:8900/ws',
        serverUrl: 'http://localhost:8900/api/v1/',
        // wsUrl: 'ws://192.168.9.160:8085/ws',
        // serverUrl: 'http://192.168.9.160:8085/api/v1/',
    }
};

const prod = {
    app: {
        version: process.env.REACT_APP_VERSION,
    },
    api: {
        wsUrl: 'ws://' + location.host + '/ws',
        serverUrl: location.origin + '/api/v1/',
    },
};

const config = process.env.REACT_APP_STAGE === 'production' ? prod : dev;

const DefaultConfig = {
    // Add common config values here
    max_attachment_size: 5000000,
    ...config
};

export default DefaultConfig;
