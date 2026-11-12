import { Component, OnInit, input, output, signal, computed } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Day, Format, FormatService, Present } from '../format.service';

@Component({
  imports: [FormsModule],
  selector: 'app-session',
  styleUrl: './session.scss',
  templateUrl: './session.html',
})
export class Session implements OnInit {
  formatId = input.required<string>();
  back = output<void>();

  format = signal<Format | null>(null);
  days = signal<Day[]>([]);
  presents = signal<Present[]>([]);

  presentsByDay = computed(() => {
    const map = new Map<string, Present[]>();
    for (const p of this.presents()) {
      const list = map.get(p.day_id) ?? [];
      list.push(p);
      map.set(p.day_id, list);
    }
    return map;
  });

  constructor(private formatService: FormatService) {}

  error = signal<string | null>(null);

  ngOnInit() {
    const id = this.formatId();
    this.formatService.get(id).subscribe({
      next: (f) => this.format.set(f),
      error: (err) => this.error.set(`Could not load session (${err.status})`),
    });
    this.formatService.getDays(id).subscribe((d) => this.days.set(d ?? []));
    this.formatService.getPresents(id).subscribe((p) => this.presents.set(p ?? []));
  }

  save() {
    const f = this.format();
    if (f) this.formatService.update(f).subscribe();
  }

  add_day() {
    const n = this.days().length;
    const d: Day = { id: crypto.randomUUID(), format_id: this.formatId(), which: n + 1, index: n };
    this.days.update((l) => [...l, d]);
    this.formatService.addDay(d).subscribe();
  }

  add_present(day_id: string) {
    const p: Present = { id: crypto.randomUUID(), day_id, present: '' };
    this.presents.update((l) => [...l, p]);
    this.formatService.addPresent(p).subscribe();
  }

  update_present(id: string, text: string) {
    this.presents.update((l) => l.map((p) => (p.id === id ? { ...p, present: text } : p)));
    this.formatService.updatePresent(id, text).subscribe();
  }

  debug = signal('');

  get() {
    const id = this.formatId();
    this.formatService.getDays(id).subscribe((days) => {
      this.formatService.getPresents(id).subscribe((presents) => {
        this.debug.set(JSON.stringify({ days, presents }, null, 2));
      });
    });
  }
}
