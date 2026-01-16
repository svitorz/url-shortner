import API from "./API";

export default {
  ping() {
    return API().get('/ping')
  },

  login(payload) {
    return API().post('/login', payload);
  },

  register(payload) {
    return API().post('/user', payload)
  },

  createLink(payload) {
    return API().post('http://localhost:8000/api/links', payload);
  }
}
