import { Component, OnInit } from '@angular/core';
import {HttpService} from "../../services/http.service";
import {ActivatedRoute, Router} from "@angular/router";

@Component({
  selector: 'app-card',
  templateUrl: './card.component.html',
  styleUrls: ['./card.component.scss']
})
export class CardComponent implements OnInit {

  constructor(
      private http: HttpService,
      private activeRoute: ActivatedRoute
  ) {
    this.activeRoute.params.subscribe((param)=> {
      this.id = param.id;
    })
  }

  id: number;
  payments: any = [{}];
  bill: any = [{}];
  visibility: boolean = true;

  ngOnInit() {
    this.getBill().then(() => {
      console.log(this.bill)
      this.bill.prenumber = this.bill.number % 10000
      switch (this.bill.type) {
        case 1: {
          this.bill.img = "assets/icons/card.png";
          this.bill.name = "Лицевой счет"
          break
        }
        case 2: {
          this.bill.img = "assets/icons/mastercard.png"
          this.bill.name = "Mastercard"
          break
        }
        case 3: {
          this.bill.img = "assets/icons/visa.png"
          this.bill.name = "Visa"
          break
        }
        case 4: {
          this.bill.img = "assets/icons/mir.png"
          this.bill.name = "МИР"
          break
        }
      }
      this.getPayments().then(() => {
        console.log(this.payments)
      }).finally(() => {
        this.visibility = false;
      })
    })
  }

  async getPayments() {
    this.payments = await this.http.getUserPaymentsByBill(this.id)
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

  async getBill() {
    this.bill = await this.http.getBillByID(this.id);
    this.bill = this.bill["body"];
  }

}
