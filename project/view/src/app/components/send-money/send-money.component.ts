import {Component, DoCheck, OnInit} from '@angular/core';
import {HttpService} from "../../services/http.service";
import {FormControl, FormGroup, Validators} from "@angular/forms";

@Component({
  selector: 'app-send-money',
  templateUrl: './send-money.component.html',
  styleUrls: ['./send-money.component.scss']
})
export class SendMoneyComponent implements OnInit, DoCheck {

  constructor(
      private http: HttpService
  ) { }

  visibility: boolean = true;
  Bills: any = [{}];
  sendMoney: FormGroup;
  isMobile: boolean = false;

  ngOnInit(): void {
    this.getBills().then().finally(() => {
      this.visibility = false;
    });

    this.sendMoney = new FormGroup({
      billID: new FormControl('', [Validators.required, Validators.required]),
      number: new FormControl('', [Validators.required, Validators.pattern(RegExp("^[0-9]{16}$"))]),
      amount: new FormControl('', [Validators.required, Validators.pattern(RegExp("^[0-9]+$"))])
    })
  }

  ngDoCheck() {
    if (document.documentElement.clientHeight > document.documentElement.clientWidth) {
      this.isMobile = true
    }
  }

  send() {

  }

  async getBills() {
    this.Bills = await this.http.getUserBills();
    this.Bills = this.Bills["body"]
    console.log(this.Bills)

    for (let item of this.Bills) {

      item.number = item.number % 10000;

      if (item.type === 1) {
          item.img = "assets/icons/card.png";
          item.name = "Лицевой счет"
      } else if (item.type === 2) {
          item.img = "assets/icons/mastercard.png"
          item.name = "Mastercard"
      } else if (item.type === 3) {
          item.img = "assets/icons/visa.png"
          item.name = "Visa"
      } else if (item.type === 4) {
          item.img = "assets/icons/mir.png"
          item.name = "МИР"
      }
    }

  }

}
