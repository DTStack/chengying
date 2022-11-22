export const imageStoreInfo = {
    address: 'http://docker.io',
    alias: 'test',
    email: 'emm@qq.com',
    id: 16,
    name: 'yuwan_test',
    password: 'admin123',
    username: 'admin'
}

export default {
    updateImageStore: {
        msg: 'ok',
        code: 0,
        data: {
            address: 'http://docker.io',
            alias: 'test',
            clusterId: 85,
            email: 'emm@qq.com',
            id: 16,
            is_default: 1,
            name: 'yuwan_test',
            password: 'admin123',
            username: 'admin'
        }
    },
    createImageStore: {
        msg: 'ok',
        code: 0,
        data: {
            address: 'http://docker.io',
            alias: 'test',
            clusterId: 85,
            email: 'emm@qq.com',
            id: 16,
            name: 'yuwan_test',
            password: 'admin123',
            username: 'admin'
        }
    },
    getImageStoreList: {
        msg: 'ok',
        code: 0,
        data: {
            count: 2,
            list: [{ address: 'http://docker.io', alias: 'test', clusterId: 85, email: 'emm@qq.com', id: 16, is_default: 1, name: 'yuwan_test', password: 'admin123', username: 'admin' }, { address: 'http://docker.io', alias: 'test1', clusterId: 85, email: '', id: 17, is_default: 0, name: 'yuwan_test1', password: 'DT#passw0rd2019', username: 'admin@dtstack.com' }]
        }
    },
    getImageStoreInfo: {
        msg: 'ok',
        code: 0,
        data: imageStoreInfo
    },
    setDefaultStore: {
        msg: 'ok',
        code: 0,
        data: { message: 'success' }
    }
}
