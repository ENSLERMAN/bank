import { Component, OnInit } from '@angular/core';
import {HttpService} from "../../services/http.service";

@Component({
  selector: 'app-payment-history',
  templateUrl: './payment-history.component.html',
  styleUrls: ['./payment-history.component.scss']
})
export class PaymentHistoryComponent implements OnInit {

  constructor(
      private http: HttpService
  ) { }

  payments: any = [{}];
  Bills: any = [{}];

  ngOnInit() {
    this.getPayments().then(() => {
      console.log(this.payments)
    })
  }

  async getPayments() {
    this.payments = await this.http.getUserPayments()
    await this.http.getUserBills().then((res) => {
      this.Bills = res["body"]
      console.log(this.Bills)
    });
    this.payments = this.payments["body"]
    for (let item of this.payments) {
      for (let i of this.Bills) {
        if (item.sender === i.number && item.recipient !== i.number) {
          item.type = "assets/icons/down.png"
        } else if (item.sender !== i.number && item.recipient === i.number) {
          item.type = "assets/icons/up.png"
        } else if (item.sender === i.number && item.recipient === i.number) {
          item.type = "assets/icons/exchange.png"
        }
      }
      item.recipient = item.recipient % 10000
      item.sender = item.sender % 10000
    }
    for (let i = 0; i < this.payments.length - 1; i++) {
      for (let j = 1; j < this.payments.length; j++) {
        if (this.payments[i].id === this.payments[j].id) {
          this.payments.splice(i, 1)
        }
      }
    }

  }

}
