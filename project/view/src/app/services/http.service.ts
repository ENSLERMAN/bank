import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from "@angular/common/http";

@Injectable({
  providedIn: 'root'
})
export class HttpService {

  headers: HttpHeaders = new HttpHeaders({
    'Content-Type': 'application/json',
    'Access-Control-Allow-Origin': '*'
  });
  options = {
    headers: this.headers,
    observe: 'response' as 'body',
    withCredentials: true
  };

  private baseURL: string = "http://localhost:8081";

  constructor(private http: HttpClient) { }

  authUser(data) {
    return this.http.post(`${this.baseURL}/session`, data, this.options).toPromise()
  }

  getUserInfo() {
    return this.http.get(`${this.baseURL}/private/whoami`).toPromise()
  }

}
