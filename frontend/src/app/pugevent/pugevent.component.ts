import { Component, OnInit } from '@angular/core';
import {ActivatedRoute} from "@angular/router";
import {PugeventService} from "../services/pugevent.service";
import {Observable} from "rxjs/Observable";
import {Pugevent} from "../shared/models/pugevent";

@Component({
  selector: 'app-pugevent',
  templateUrl: './pugevent.component.html',
  styleUrls: ['./pugevent.component.scss']
})
export class PugeventComponent implements OnInit {
  public eventsObservable: Observable<Pugevent>;

  constructor(
    private route: ActivatedRoute,
    private pugeventService: PugeventService
  ) { }

  ngOnInit() {
    this.eventsObservable = this.route.params.switchMap(() => {
      return this.pugeventService.list();
    });
  }

}
