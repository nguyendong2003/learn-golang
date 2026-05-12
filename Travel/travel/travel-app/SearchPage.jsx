// SearchPage — Filter sidebar + Tour/Hotel/Combo results
function SearchPage({ navigate, params }) {
  const type0 = params?.type || 'tour';
  const [activeType, setActiveType] = useState(type0);
  const [filters, setFilters] = useState({ dest: params?.dest || '', priceMax: 10000000, duration: [], rating: 0, sort: 'popular' });
  const [saved, setSaved] = useState({});

  const setF = (k, v) => setFilters(f => ({ ...f, [k]: v }));
  const toggleDur = (d) => setF('duration', filters.duration.includes(d) ? filters.duration.filter(x => x !== d) : [...filters.duration, d]);

  const items = activeType === 'hotel' ? HOTELS : activeType === 'combo' ? COMBOS : TOURS;
  const filtered = items.filter(i => {
    if (filters.dest && !((i.dest||i.loc||i.name||'').toLowerCase().includes(filters.dest.toLowerCase()))) return false;
    if (i.price > filters.priceMax) return false;
    if (filters.rating > 0 && i.rating < filters.rating) return false;
    return true;
  });

  const typeTabs = [
    { id:'tour', label:'Tour', count: TOURS.length },
    { id:'hotel', label:'Khách sạn', count: HOTELS.length },
    { id:'combo', label:'Combo', count: COMBOS.length },
  ];

  const durations = ['1-2 ngày','3-4 ngày','5-7 ngày','Trên 7 ngày'];
  const sorts = [
    { id:'popular', label:'Phổ biến nhất' },
    { id:'price-asc', label:'Giá thấp → cao' },
    { id:'price-desc', label:'Giá cao → thấp' },
    { id:'rating', label:'Đánh giá cao nhất' },
    { id:'newest', label:'Mới nhất' },
  ];

  return (
    <div className="page-anim" style={{ paddingTop:68 }}>
      {/* ── Top search bar ── */}
      <div style={{ background:'var(--primary)', padding:'20px 0' }}>
        <div className="container">
          <div style={{ display:'flex', alignItems:'center', gap:12, background:'rgba(255,255,255,.12)', borderRadius:'var(--r-lg)', padding:'10px 16px' }}>
            <span style={{ color:'rgba(255,255,255,.7)' }}><Ico.search /></span>
            <input defaultValue={filters.dest || ''} onChange={e => setF('dest', e.target.value)}
              placeholder="Tìm tour, khách sạn, combo..." 
              style={{ flex:1, background:'transparent', border:'none', outline:'none', color:'#fff', fontSize:15, fontWeight:500 }} />
            <div style={{ display:'flex', gap:8 }}>
              {typeTabs.map(t => (
                <button key={t.id} onClick={() => setActiveType(t.id)}
                  style={{ padding:'6px 14px', borderRadius:'var(--r-full)', border:'none', cursor:'pointer', fontSize:13, fontWeight:700, transition:'all .15s',
                    background: activeType===t.id ? '#fff' : 'rgba(255,255,255,.15)',
                    color: activeType===t.id ? 'var(--primary)' : 'rgba(255,255,255,.85)',
                  }}>{t.label} <span style={{ opacity:.7 }}>({t.count})</span></button>
              ))}
            </div>
          </div>
        </div>
      </div>

      <div className="container" style={{ paddingTop:28, paddingBottom:40 }}>
        <div style={{ display:'flex', gap:28, alignItems:'flex-start' }}>

          {/* ── Filter Sidebar ── */}
          <aside style={{ width:256, flexShrink:0, position:'sticky', top:90 }}>
            <div className="card-flat" style={{ padding:'20px', borderRadius:'var(--r-lg)' }}>
              <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center', marginBottom:18 }}>
                <span style={{ fontWeight:700, fontSize:15, display:'flex', alignItems:'center', gap:6 }}><Ico.filter /> Bộ lọc</span>
                <button onClick={() => setFilters({ dest:'', priceMax:10000000, duration:[], rating:0, sort:'popular' })}
                  style={{ fontSize:12, color:'var(--primary)', background:'none', border:'none', cursor:'pointer', fontWeight:600 }}>Xóa tất cả</button>
              </div>

              {/* Điểm đến */}
              <FilterSection title="Điểm đến">
                <input className="input" placeholder="Nhập điểm đến..." value={filters.dest} onChange={e => setF('dest', e.target.value)} style={{ fontSize:13 }} />
                <div style={{ display:'flex', gap:6, flexWrap:'wrap', marginTop:8 }}>
                  {['Hạ Long','Sapa','Phú Quốc','Đà Nẵng','Đà Lạt'].map(d => (
                    <button key={d} onClick={() => setF('dest', filters.dest===d ? '' : d)}
                      className={'chip' + (filters.dest===d ? ' active' : '')} style={{ fontSize:12, padding:'3px 10px' }}>{d}</button>
                  ))}
                </div>
              </FilterSection>

              {/* Khoảng giá */}
              <FilterSection title="Khoảng giá">
                <div style={{ marginBottom:6 }}>
                  <div style={{ display:'flex', justifyContent:'space-between', fontSize:12, color:'var(--text2)', marginBottom:4 }}>
                    <span>0đ</span>
                    <span style={{ fontWeight:700, color:'var(--primary)' }}>{formatPrice(filters.priceMax)}</span>
                  </div>
                  <input type="range" min={500000} max={20000000} step={500000} value={filters.priceMax} onChange={e => setF('priceMax', +e.target.value)}
                    style={{ width:'100%', accentColor:'var(--primary)' }} />
                </div>
                {[[0,2000000,'Dưới 2 triệu'],[2000000,5000000,'2 — 5 triệu'],[5000000,10000000,'5 — 10 triệu'],[10000000,99999999,'Trên 10 triệu']].map(([min,max,label]) => (
                  <label key={label} style={{ display:'flex', gap:8, alignItems:'center', fontSize:13, cursor:'pointer', padding:'3px 0' }}>
                    <input type="radio" name="price" onChange={() => setF('priceMax', max)} style={{ accentColor:'var(--primary)' }} /> {label}
                  </label>
                ))}
              </FilterSection>

              {/* Thời lượng (tour only) */}
              {activeType === 'tour' && (
                <FilterSection title="Thời lượng">
                  {durations.map(d => (
                    <label key={d} style={{ display:'flex', gap:8, alignItems:'center', fontSize:13, cursor:'pointer', padding:'3px 0' }}>
                      <input type="checkbox" checked={filters.duration.includes(d)} onChange={() => toggleDur(d)} style={{ accentColor:'var(--primary)' }} /> {d}
                    </label>
                  ))}
                </FilterSection>
              )}

              {/* Đánh giá */}
              <FilterSection title="Xếp hạng tối thiểu" last>
                {[5,4,3,0].map(r => (
                  <label key={r} style={{ display:'flex', gap:8, alignItems:'center', fontSize:13, cursor:'pointer', padding:'3px 0' }}>
                    <input type="radio" name="rating" checked={filters.rating===r} onChange={() => setF('rating', r)} style={{ accentColor:'var(--primary)' }} />
                    {r > 0 ? <><Stars score={r} size={12} /> trở lên</> : 'Tất cả'}
                  </label>
                ))}
              </FilterSection>
            </div>
          </aside>

          {/* ── Results ── */}
          <div style={{ flex:1, minWidth:0 }}>
            {/* Sort + count bar */}
            <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center', marginBottom:20 }}>
              <div>
                <span style={{ fontWeight:700, fontSize:16 }}>{filtered.length}</span>
                <span style={{ color:'var(--text2)', fontSize:14, marginLeft:6 }}>kết quả{filters.dest ? ' cho "' + filters.dest + '"' : ''}</span>
              </div>
              <div style={{ display:'flex', alignItems:'center', gap:8 }}>
                <span style={{ fontSize:13, color:'var(--text2)' }}>Sắp xếp:</span>
                <select value={filters.sort} onChange={e => setF('sort', e.target.value)}
                  style={{ border:'1.5px solid var(--border)', borderRadius:'var(--r-sm)', padding:'6px 12px', fontSize:13, fontWeight:600, color:'var(--text)', background:'#fff', cursor:'pointer', outline:'none' }}>
                  {sorts.map(s => <option key={s.id} value={s.id}>{s.label}</option>)}
                </select>
              </div>
            </div>

            {/* Grid */}
            {filtered.length === 0 ? (
              <div style={{ textAlign:'center', padding:'60px 0', color:'var(--text3)' }}>
                <div style={{ fontSize:48, marginBottom:12 }}>🔍</div>
                <div style={{ fontSize:16, fontWeight:600 }}>Không tìm thấy kết quả</div>
                <div style={{ fontSize:14, marginTop:4 }}>Thử thay đổi bộ lọc hoặc tìm kiếm khác</div>
              </div>
            ) : (
              <div style={{ display:'grid', gridTemplateColumns:'repeat(3, 1fr)', gap:20 }}>
                {activeType === 'tour' && filtered.map(t => (
                  <TourCard key={t.id} t={t} onClick={() => navigate('tour-detail', { tour:t })} />
                ))}
                {activeType === 'hotel' && filtered.map(h => (
                  <HotelCard key={h.id} h={h} onClick={() => navigate('hotel-detail', { hotel:h })} />
                ))}
                {activeType === 'combo' && filtered.map(c => (
                  <SearchComboCard key={c.id} c={c} onClick={() => navigate('combo-detail', { combo:c })} />
                ))}
              </div>
            )}

            {/* Pagination */}
            {filtered.length > 0 && (
              <div style={{ display:'flex', justifyContent:'center', gap:6, marginTop:32 }}>
                {[1,2,3,'...',8].map((p,i) => (
                  <button key={i}
                    style={{ width:36, height:36, borderRadius:'var(--r-sm)', border: p===1 ? 'none' : '1.5px solid var(--border)', cursor:'pointer', fontWeight: p===1 ? 700 : 400, fontSize:14,
                      background: p===1 ? 'var(--primary)' : 'var(--surface)',
                      color: p===1 ? '#fff' : 'var(--text2)' }}>{p}</button>
                ))}
                <button style={{ height:36, padding:'0 14px', borderRadius:'var(--r-sm)', border:'1.5px solid var(--border)', cursor:'pointer', fontSize:13, background:'#fff', color:'var(--text2)', display:'flex', alignItems:'center', gap:4 }}>Tiếp <Ico.chevR /></button>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

function FilterSection({ title, children, last }) {
  return (
    <div style={{ borderBottom: last ? 'none' : '1px solid var(--border-lt)', paddingBottom: last ? 0 : 16, marginBottom: last ? 0 : 16 }}>
      <div style={{ fontWeight:700, fontSize:13, color:'var(--text)', marginBottom:10, letterSpacing:.2 }}>{title}</div>
      <div style={{ display:'flex', flexDirection:'column', gap:0 }}>{children}</div>
    </div>
  );
}

function SearchComboCard({ c, onClick }) {
  return (
    <div className="card" style={{ cursor:'pointer' }} onClick={onClick}>
      <Img bg={c.bg} stripe={c.stripe} label={c.label} style={{ height:160, width:'100%' }} />
      <div style={{ padding:'14px 16px' }}>
        <span className="badge badge-danger" style={{ marginBottom:8 }}>Tiết kiệm {c.disc}%</span>
        <div style={{ fontWeight:700, fontSize:14, lineHeight:1.4, marginBottom:8 }}>{c.name}</div>
        <div style={{ fontSize:12, color:'var(--text3)', marginBottom:10, display:'flex', gap:4, alignItems:'center' }}>
          <Ico.plane /> {c.from} → {c.dest} · {c.duration}
        </div>
        <div style={{ display:'flex', alignItems:'center', justifyContent:'space-between' }}>
          <div>
            <div style={{ fontSize:12, color:'var(--text3)', textDecoration:'line-through' }}>{formatPrice(c.oldPrice)}</div>
            <div style={{ fontSize:16, fontWeight:800, color:'var(--accent)' }}>{formatPrice(c.price)}</div>
          </div>
          <button className="btn btn-primary btn-sm">Xem combo</button>
        </div>
      </div>
    </div>
  );
}

Object.assign(window, { SearchPage });
