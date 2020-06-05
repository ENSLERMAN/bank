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
  visibility: boolean = true;

  ngOnInit() {
    this.getPayments().then(() => {
      console.log(this.payments)
    }).finally(() => {
      this.visibility = false;
    })
  }

  async getPayments() {
    this.payments = await this.http.getUserPayments()
    this.payments = this.payments["body"]

    for (let item of this.payments) {
      if (item.type === 1) {
        item.img = "assets/icons/down.png"
        item.text = "Списание"
      } else if (item.type === 2) {
        item.img = "assets/icons/up.png"
        item.text = "Пополнение"
      } else if (item.type === 3) {
        item.img = "assets/icons/exchange.png"
        item.text = "Перевод между счетами"
      }
      item.recipient = item.recipient % 10000
      item.sender = item.sender % 10000
    }

  }



}
