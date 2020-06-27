import {Component, DoCheck, OnInit} from '@angular/core';
import {HttpService} from "../../services/http.service";
import {Router} from "@angular/router";
import {MatSnackBar} from "@angular/material/snack-bar";

@Component({
  selector: 'app-payments',
  templateUrl: './payments.component.html',
  styleUrls: ['./payments.component.scss']
})
export class PaymentsComponent implements OnInit, DoCheck {

  constructor(
      private http: HttpService,
      private router: Router,
      private _snackBar: MatSnackBar
  ) { }

  bills: any;
  bill: any;
  amount: number;
  isMobile: boolean;
  visibility: boolean = false;

  ngOnInit(): void {
    if (document.documentElement.clientHeight > document.documentElement.clientWidth) {
      this.isMobile = true
    }
  }

  ngDoCheck(): void {
    if (document.documentElement.clientHeight > document.documentElement.clientWidth) {
      this.isMobile = true
    }
  }

  async onPay() {
    this.visibility = true;
    this.bills = await this.http.getUserBills()
    this.bills = this.bills["body"];
    console.log(this.bills)

    for (let item of this.bills) {
      if (item.money >= 5000) {
        this.bill = item;
        break;
      }
    }
    console.log(this.bill)
    this.amount = this.randomNumber()

    await this.http.sendMoney({
      "bill_id": this.bill.id,
      "number_dest": 1000000000001001,
      "amount": this.amount,
    }).then(
        (res) => {
          if (res["status"] == 200) {
            console.log("Transfer successful");
            this._snackBar.open(`Списание на сумму: ${this.amount} ₽`, `X`, {
              duration: 5000,
            });
          }
        },
        (err) => {
          console.log("Запрос на перевод не прошёл")
        }
    ).finally(()=> {
      this.visibility = false;
      this.router.navigate(['/main']);
    });
  }

  randomNumber() {
    let rand = 1 + Math.random() * (5000 + 1 - 1);
    return Math.floor(rand);
  }

}
