
import axios from 'axios'
import { resolve } from 'url';

var axiosi = axios.create({
    headers: {
        'Access-Control-Allow-Origin': '*',
        'Content-Type': 'application/json',
    },
    timeout: 30000,
    baseURL: 'http://127.0.0.1:3000'
})

// 请求拦截器
axiosi.interceptors.request.use(config => {
    return config
}, error => {
    console.log(error)
    return Promise.reject(error)
})
// 响应拦截器
axiosi.interceptors.response.use(response => {
    return response.data
}, error => {
    console.log('err' + error)
    return Promise.reject(error)
})
export default axiosi;
/**
 * post 请求方法
 * @param {*} url 
 * @param {*} data 
 */
export function post(url, data = {}) {
    return new Promise((resplve, reject) => {
        axiosi.post(url, data)
            .then(response => {
                resolve(response.data)
            }, err => {
                reject(err)
            })
    })
}

/**
 * get 请求方法
 * @param {*} url 
 * @param {*} data 
 */
export function get(url, data = {}) {
    return new Promise((resolve, reject) => {
        axiosi.get(url, {
            params: data
        })
            .then(response => {
                resolve(response)
            })
            .catch(err => {
                reject(err)
            })
    })
}

/**
 * 封装所有请求
 * @param {*} method 
 * @param {*} url 
 * @param {*} data 
 * @param {*} headers 
 * @returns {Promise}
 */
export function request(method, url, data = {}, headers) {
    return new Promise((resolve, reject) => {
        axiosi({
            method: method || 'post',
            url: url,
            data: method === 'get' ? { params: data } : data,
            headers: headers || { 'Content-Type': 'application/json' },
        })
            .then(response => {
                resolve(response.data);
            })
            .catch(err => {
                reject(err)
            })
    })
}

