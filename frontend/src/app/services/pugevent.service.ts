import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/map';

import {Pugevent} from '../shared/models/pugevent';

@Injectable()
export class PugeventService {
  private API = '/api/1';
  // private API = 'https://gcpug.jp/api/1';

  constructor(public http: HttpClient) {}

  list(): Observable<Pugevent> {
    return this.http.get<Pugevent>(`${this.API}/event`, {});
  }
}
