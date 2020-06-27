import {Component, DoCheck, OnInit} from '@angular/core';
import {HttpService} from "../../services/http.service";
import {FormControl, FormGroup, Validators} from "@angular/forms";
import {ActivatedRoute, Router} from "@angular/router";

@Component({
  selector: 'app-send-money',
  templateUrl: './send-money.component.html',
  styleUrls: ['./send-money.component.scss']
})
export class SendMoneyComponent implements OnInit, DoCheck {

  constructor(
      private http: HttpService,
      private router: Router,
      private activeRoute: ActivatedRoute
  ) {
      activeRoute.queryParams.subscribe(
          (queryParam: any) => {
              this.selected = queryParam['selected']
          }
      )
  }

  visibility: boolean = true;
  Bills: any = [{}];
  sendMoney: FormGroup;
  sendMoney2: FormGroup;
  isMobile: boolean = false;
  selected: any;
  restBills: any;

  ngOnInit(): void {
    this.getBills().then(()=> {
        this.getRestOfTheBills().then()
    }).finally(() => {
      this.visibility = false;
    });

    this.sendMoney = new FormGroup({
      billID: new FormControl('', [Validators.required, Validators.required]),
      number: new FormControl('', [Validators.required, Validators.pattern(RegExp("^[0-9]{16}$"))]),
      amount: new FormControl('', [Validators.required, Validators.pattern(RegExp("^[0-9]+$"))])
    });

    this.sendMoney2 = new FormGroup({
        billIDsend: new FormControl('', [Validators.required, Validators.required]),
        billIDrec: new FormControl('', [Validators.required, Validators.required]),
        amount: new FormControl('', [Validators.required, Validators.pattern(RegExp("^[0-9]+$"))])
    });
  }

  ngDoCheck() {
    if (document.documentElement.clientHeight > document.documentElement.clientWidth) {
      this.isMobile = true
    }
  }

  async getRestOfTheBills() {
      this.restBills = await this.http.getRestOfTheBills();
      this.restBills = this.restBills["body"]
      console.log(this.restBills)
  }

  async send() {
    await this.http.sendMoney({
      "bill_id": this.sendMoney.value.billID,
      "number_dest": parseInt(this.sendMoney.value.number),
      "amount": parseInt(this.sendMoney.value.amount),
    }).then(
        (res) => {
          if (res["status"] == 200) {
            console.log("Transfer successful");
            this.router.navigate(['/main']);
          }
        },
        (err) => {
          console.log("Запрос на перевод не прошёл")
        }
    );
  }

    async send2() {
        await this.http.sendMoney({
            "bill_id": this.sendMoney2.value.billIDsend,
            "number_dest": this.sendMoney2.value.billIDrec,
            "amount": parseInt(this.sendMoney2.value.amount),
        }).then(
            (res) => {
                if (res["status"] == 200) {
                    console.log("Transfer successful");
                    this.router.navigate(['/main']);
                }
            },
            (err) => {
                console.log("Запрос на перевод не прошёл")
            }
        );
    }

  async getBills() {
    this.Bills = await this.http.getUserBills();
    this.Bills = this.Bills["body"]
    console.log(this.Bills)

    for (let item of this.Bills) {

      item.prenumber = item.number % 10000;

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
