import http from '@/http'
import { Api } from '../common/enum'


export interface IConfig {
    api_key: string
    base_url: string
}

export const getConfigApi = () => {
    return http.get(`${Api.maxkbconfig}`)
};

export const updateConfigApi = (params: IConfig) => {
    return http.patch(`${Api.maxkbconfig}`, params)
}
