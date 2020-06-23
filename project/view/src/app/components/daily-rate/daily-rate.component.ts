import { Component, OnInit } from '@angular/core';
import {HttpService} from "../../services/http.service";
import {HttpClient} from "@angular/common/http";

@Component({
  selector: 'app-daily-rate',
  templateUrl: './daily-rate.component.html',
  styleUrls: ['./daily-rate.component.scss']
})
export class DailyRateComponent implements OnInit {

  constructor(
      private http: HttpClient
  ) { }

  visibility: boolean = true;
  rates: any = {};

  async ngOnInit() {
    this.visibility = true;
    setTimeout(()=>{
      this.visibility = false;
    }, 1500)
  }

}
