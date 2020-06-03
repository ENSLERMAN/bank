import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from "@angular/common/http";

@Injectable({
  providedIn: 'root'
})
export class HttpService {

  headers: HttpHeaders = new HttpHeaders({
    'Content-Type': 'application/json'
  });
  options = {
    headers: this.headers,
    observe: 'response' as 'body',
    withCredentials: true
  };

  private baseURL: string = "http://localhost:8081";

  constructor(private http: HttpClient) { }

  authUser(data) {
    return this.http.post(`${this.baseURL}/sessions`, data, this.options).toPromise();
  }

  getUserInfo() {
    return this.http.get(`${this.baseURL}/private/whoami`).toPromise()
  }

  getUserBills() {
    return this.http.get(`${this.baseURL}/private/get_bills`, this.options).toPromise();
  }

  createBill(data) {
    return this.http.post(`${this.baseURL}/private/create_bill`, data, this.options).toPromise();
  }

  getUserPayments() {
    return this.http.get(`${this.baseURL}/private/get_payments`, this.options).toPromise();
  }

}
