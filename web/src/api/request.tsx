import { axiosPost, axiosPut, axiosGet, axiosPatch, axiosDelete } from './axios'

export async function Get(url: string) {
    if (isDev()) {
        
    }
    return await axiosGet(url)
}

function isDev(): boolean {
    return process.env.NODE_ENV === 'development'
}