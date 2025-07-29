import JSEncrypt from 'jsencrypt'
import CryptoJS from 'crypto-js'
import Cookies from 'js-cookie'

export function encryptPassword(password: string): string {
  const publicKeyBase64 = Cookies.get('mkics_public_key')
  if (!publicKeyBase64) return password

  const publicKey = atob(publicKeyBase64.replaceAll('"', ''))

  const aesKey = CryptoJS.lib.WordArray.random(16)
  const aesKeyBase64 = CryptoJS.enc.Base64.stringify(aesKey)

  const iv = CryptoJS.lib.WordArray.random(16)
  const ivBase64 = CryptoJS.enc.Base64.stringify(iv)

  const encrypted = CryptoJS.AES.encrypt(password, aesKey, {
    iv,
    mode: CryptoJS.mode.CBC,
    padding: CryptoJS.pad.Pkcs7
  })

  const cipherText = encrypted.toString()

  const jsEncrypt = new JSEncrypt()
  jsEncrypt.setPublicKey(publicKey)
  const encryptedKey = jsEncrypt.encrypt(aesKeyBase64)

  return `${encryptedKey}:${ivBase64}:${cipherText}`
}
