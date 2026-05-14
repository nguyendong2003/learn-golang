// HomePage — Hero, Destinations, Tours, Flash Sale, Hotels, Combos, Blog
function HomePage({ navigate }) {
  const [tab, setTab] = useState('tour');
  const [dest, setDest] = useState('');
  const [dateIn, setDateIn] = useState('');
  const [dateOut, setDateOut] = useState('');
  const [adults, setAdults] = useState(2);

  function doSearch() { navigate('search', { type: tab, dest }); }

  return (
    <div className="page-anim">
      {/* ── HERO ── */}
      <section style={{ position:'relative', height:'92vh', minHeight:620, display:'flex', alignItems:'flex-end', paddingBottom:64 }}>
        <Img bg="#0c3d5a" stripe="#082d43" label="" style={{ position:'absolute', inset:0, width:'100%', height:'100%' }} />
        {/* Gradient overlay */}
        <div style={{ position:'absolute', inset:0, background:'linear-gradient(180deg, rgba(0,0,0,.15) 0%, rgba(0,0,0,.05) 40%, rgba(0,0,0,.72) 100%)' }} />
        <div className="container" style={{ position:'relative', zIndex:1, width:'100%' }}>
          {/* Headline */}
          <div style={{ marginBottom:32, maxWidth:680 }}>
            <div style={{ display:'flex', alignItems:'center', gap:10, marginBottom:14 }}>
              <span style={{ background:'var(--accent)', color:'#fff', padding:'4px 12px', borderRadius:'var(--r-full)', fontSize:12, fontWeight:700, letterSpacing:.5, whiteSpace:'nowrap' }}>DU LỊCH VIỆT NAM</span>
            </div>
            <h1 style={{ fontSize:52, fontWeight:900, color:'#fff', lineHeight:1.1, letterSpacing:'-0.03em', textWrap:'balance' }}>
              Khám phá Việt Nam<br/><span style={{ color:'rgba(255,255,255,0.78)' }}>tươi đẹp & trọn vẹn</span>
            </h1>
            <p style={{ marginTop:16, fontSize:18, color:'rgba(255,255,255,0.78)', fontWeight:400, maxWidth:480 }}>
              Hàng nghìn tour, khách sạn và combo du lịch hấp dẫn — chỉ trong một nơi.
            </p>
          </div>

          {/* Search widget */}
          <SearchWidget tab={tab} setTab={setTab} dest={dest} setDest={setDest}
            dateIn={dateIn} setDateIn={setDateIn} dateOut={dateOut} setDateOut={setDateOut}
            adults={adults} setAdults={setAdults} onSearch={doSearch} />

          {/* Quick destination chips */}
          <div style={{ display:'flex', gap:8, marginTop:18, flexWrap:'wrap' }}>
            {['Hạ Long','Sapa','Phú Quốc','Đà Nẵng','Đà Lạt','Nha Trang'].map(d => (
              <button key={d} onClick={() => { setDest(d); doSearch(); }}
                style={{ padding:'5px 14px', borderRadius:'var(--r-full)', border:'1.5px solid rgba(255,255,255,.45)', background:'rgba(255,255,255,.12)', color:'#fff', fontSize:13, fontWeight:500, cursor:'pointer', backdropFilter:'blur(6px)', transition:'all .15s' }}
                onMouseEnter={e => { e.target.style.background='rgba(255,255,255,.22)'; e.target.style.borderColor='rgba(255,255,255,.7)'; }}
                onMouseLeave={e => { e.target.style.background='rgba(255,255,255,.12)'; e.target.style.borderColor='rgba(255,255,255,.45)'; }}
              >{d}</button>
            ))}
          </div>
        </div>
      </section>

      {/* ── ĐIỂM ĐẾN NỔI BẬT ── */}
      <section className="section-sm">
        <div className="container">
          <div style={{ display:'flex', justifyContent:'space-between', alignItems:'flex-end', marginBottom:24 }}>
            <div>
              <h2 className="section-title">Điểm đến nổi bật</h2>
              <p className="section-sub">Những vùng đất được yêu thích nhất tại Việt Nam</p>
            </div>
            <button onClick={() => navigate('search',{type:'tour'})} className="btn btn-ghost btn-sm" style={{ gap:4 }}>Xem thêm <Ico.chevR /></button>
          </div>
          <div className="h-scroll">
            {DESTINATIONS.map(d => (
              <DestCard key={d.id} d={d} onClick={() => navigate('search',{type:'tour',dest:d.name})} />
            ))}
          </div>
        </div>
      </section>

      {/* ── TOUR HOT ── */}
      <section className="section" style={{ background:'var(--surface)', borderTop:'1px solid var(--border-lt)' }}>
        <div className="container">
          <div style={{ display:'flex', justifyContent:'space-between', alignItems:'flex-end', marginBottom:24 }}>
            <div>
              <h2 className="section-title">Tour hot trong tuần</h2>
              <p className="section-sub">Được đặt nhiều nhất — đừng bỏ lỡ chỗ ngồi</p>
            </div>
            <button onClick={() => navigate('search',{type:'tour'})} className="btn btn-outline btn-sm">Xem tất cả tour</button>
          </div>
          <div className="h-scroll" style={{ gap:20 }}>
            {TOURS.slice(0,6).map(t => <TourCard key={t.id} t={t} onClick={() => navigate('tour-detail',{tour:t})} />)}
          </div>
        </div>
      </section>

      {/* ── FLASH SALE ── */}
      <FlashSaleBanner navigate={navigate} />

      {/* ── KHÁCH SẠN ── */}
      <section className="section">
        <div className="container">
          <div style={{ display:'flex', justifyContent:'space-between', alignItems:'flex-end', marginBottom:24 }}>
            <div>
              <h2 className="section-title">Khách sạn đề xuất</h2>
              <p className="section-sub">Chỗ nghỉ lý tưởng với giá tốt nhất được đảm bảo</p>
            </div>
            <button onClick={() => navigate('search',{type:'hotel'})} className="btn btn-outline btn-sm">Xem tất cả KS</button>
          </div>
          <div className="grid-4">
            {HOTELS.slice(0,4).map(h => <HotelCard key={h.id} h={h} onClick={() => navigate('hotel-detail',{hotel:h})} />)}
          </div>
        </div>
      </section>

      {/* ── COMBO ── */}
      <section className="section-sm" style={{ background:'var(--primary-xlt)', borderTop:'1px solid var(--border-lt)', borderBottom:'1px solid var(--border-lt)' }}>
        <div className="container">
          <div style={{ display:'flex', justifyContent:'space-between', alignItems:'flex-end', marginBottom:24 }}>
            <div>
              <h2 className="section-title">Combo trọn gói tiết kiệm</h2>
              <p className="section-sub">Bay + Ở + Tour — đặt một lần, tiết kiệm hơn</p>
            </div>
            <button onClick={() => navigate('search',{type:'combo'})} className="btn btn-outline btn-sm">Xem tất cả combo</button>
          </div>
          <div className="grid-3">
            {COMBOS.map(c => <ComboCard key={c.id} c={c} navigate={navigate} />)}
          </div>
        </div>
      </section>

      {/* ── BLOG ── */}
      <section className="section">
        <div className="container">
          <div style={{ display:'flex', justifyContent:'space-between', alignItems:'flex-end', marginBottom:24 }}>
            <div>
              <h2 className="section-title">Blog du lịch</h2>
              <p className="section-sub">Cẩm nang, kinh nghiệm và câu chuyện từ khắp Việt Nam</p>
            </div>
            <button className="btn btn-ghost btn-sm" style={{ gap:4 }}>Đọc thêm <Ico.chevR /></button>
          </div>
          <div className="grid-3">
            {BLOGS.map(b => <BlogCard key={b.id} b={b} />)}
          </div>
        </div>
      </section>

      {/* ── AFFILIATE CTA ── */}
      <section style={{ background:'var(--primary)', padding:'56px 0' }}>
        <div className="container" style={{ display:'flex', alignItems:'center', justifyContent:'space-between', gap:40 }}>
          <div>
            <div style={{ fontSize:12, fontWeight:700, color:'rgba(255,255,255,.6)', letterSpacing:2, marginBottom:8, textTransform:'uppercase' }}>Chương trình cộng tác viên</div>
            <h3 style={{ fontSize:30, fontWeight:800, color:'#fff', letterSpacing:'-.02em', marginBottom:10 }}>Kiếm hoa hồng khi chia sẻ link du lịch</h3>
            <p style={{ fontSize:15, color:'rgba(255,255,255,.75)', maxWidth:520 }}>Tham gia ngay chương trình Affiliate của TravelViet — chia sẻ link, nhận hoa hồng hấp dẫn cho mỗi booking thành công.</p>
          </div>
          <div style={{ flexShrink:0, display:'flex', gap:12 }}>
            <button className="btn btn-accent btn-lg">Đăng ký ngay</button>
            <button className="btn btn-outline-w btn-lg">Tìm hiểu thêm</button>
          </div>
        </div>
      </section>

      <Footer />
    </div>
  );
}

// ── Search Widget ──────────────────────────────
function SearchWidget({ tab, setTab, dest, setDest, dateIn, setDateIn, dateOut, setDateOut, adults, setAdults, onSearch }) {
  const tabs = [{ id:'tour', label:'Tour' }, { id:'hotel', label:'Khách sạn' }, { id:'combo', label:'Combo' }];
  return (
    <div style={{ background:'#fff', borderRadius:'var(--r-xl)', boxShadow:'var(--sh-xl)', overflow:'hidden', maxWidth:920 }}>
      {/* Tabs */}
      <div style={{ display:'flex', borderBottom:'1px solid var(--border-lt)' }}>
        {tabs.map(t => (
          <button key={t.id} onClick={() => setTab(t.id)}
            style={{ flex:1, padding:'14px 0', border:'none', cursor:'pointer', fontWeight:700, fontSize:14, letterSpacing:.2, transition:'all .15s',
              background: tab === t.id ? 'var(--primary)' : 'transparent',
              color: tab === t.id ? '#fff' : 'var(--text2)',
              borderBottom: tab === t.id ? '3px solid var(--primary)' : '3px solid transparent'
            }}>{t.label}</button>
        ))}
      </div>
      {/* Fields */}
      <div style={{ display:'flex', alignItems:'center', padding:'16px 20px', gap:0 }}>
        <SearchField icon={<Ico.pin />} label="Điểm đến" placeholder="Nhập điểm đến..." value={dest} onChange={e => setDest(e.target.value)} flex={2} />
        <div style={{ width:1, height:44, background:'var(--border-lt)', flexShrink:0 }} />
        <SearchField icon={<Ico.clock />} label={tab==='hotel' ? 'Nhận phòng' : 'Ngày khởi hành'} placeholder="dd/mm/yyyy" type="date" value={dateIn} onChange={e => setDateIn(e.target.value)} />
        {tab === 'hotel' && <>
          <div style={{ width:1, height:44, background:'var(--border-lt)', flexShrink:0 }} />
          <SearchField icon={<Ico.clock />} label="Trả phòng" placeholder="dd/mm/yyyy" type="date" value={dateOut} onChange={e => setDateOut(e.target.value)} />
        </>}
        <div style={{ width:1, height:44, background:'var(--border-lt)', flexShrink:0 }} />
        <SearchField icon={<Ico.person />} label={tab==='hotel' ? 'Số khách / phòng' : 'Số người'} placeholder="2 người lớn" readOnly value={adults + ' người lớn'} />
        <div style={{ flexShrink:0, paddingLeft:12 }}>
          <button onClick={onSearch} className="btn btn-accent btn-lg" style={{ paddingLeft:28, paddingRight:28, gap:8, borderRadius:'var(--r-md)' }}>
            <Ico.search /> Tìm kiếm
          </button>
        </div>
      </div>
    </div>
  );
}

function SearchField({ icon, label, placeholder, value, onChange, type, readOnly, flex = 1 }) {
  return (
    <div style={{ flex, padding:'0 20px', minWidth:0 }}>
      <div style={{ fontSize:11, fontWeight:700, color:'var(--text3)', letterSpacing:.5, marginBottom:3, textTransform:'uppercase' }}>{label}</div>
      <div style={{ display:'flex', alignItems:'center', gap:6 }}>
        <span style={{ color:'var(--primary)', flexShrink:0 }}>{icon}</span>
        <input readOnly={readOnly} type={type||'text'} value={value} onChange={onChange}
          placeholder={placeholder}
          style={{ border:'none', outline:'none', fontSize:14, fontWeight:500, color:'var(--text)', width:'100%', background:'transparent', cursor: readOnly?'pointer':'text' }} />
      </div>
    </div>
  );
}

// ── Destination Card ──
function DestCard({ d, onClick }) {
  return (
    <div onClick={onClick} style={{ cursor:'pointer', flexShrink:0, textAlign:'center' }}>
      <Img bg={d.bg} stripe={d.stripe} label={d.label}
        style={{ width:120, height:120, borderRadius:'50%', marginBottom:10, boxShadow:'var(--sh-md)', transition:'transform .2s, box-shadow .2s' }}
        className="dest-img" />
      <div style={{ fontWeight:700, fontSize:14, color:'var(--text)' }}>{d.name}</div>
      <div style={{ fontSize:12, color:'var(--text3)', marginTop:1 }}>{d.count} tour</div>
    </div>
  );
}

// ── Tour Card ──
function TourCard({ t, onClick }) {
  const [saved, setSaved] = useState(false);
  return (
    <div className="card" onClick={onClick} style={{ width:272, flexShrink:0, cursor:'pointer' }}>
      <div style={{ position:'relative' }}>
        <Img bg={t.bg} stripe={t.stripe} label={t.label} style={{ height:168, width:'100%' }} />
        {t.badge && (
          <span className={'badge badge-'+t.badgeV} style={{ position:'absolute', top:10, left:10 }}>{t.badge}</span>
        )}
        {t.seats <= 5 && t.seats > 0 && !t.badge && (
          <span className="badge badge-danger" style={{ position:'absolute', top:10, left:10 }}>Còn {t.seats} chỗ</span>
        )}
        <button onClick={e => { e.stopPropagation(); setSaved(!saved); }}
          style={{ position:'absolute', top:10, right:10, width:32, height:32, borderRadius:'50%', border:'none', background:'rgba(255,255,255,.9)', cursor:'pointer', display:'flex', alignItems:'center', justifyContent:'center', color: saved?'var(--danger)':'var(--text3)' }}>
          <Ico.heart filled={saved} />
        </button>
        <div style={{ position:'absolute', bottom:0, left:0, right:0, height:50, background:'linear-gradient(0deg, rgba(0,0,0,.5), transparent)' }} />
        <div style={{ position:'absolute', bottom:8, left:10, display:'flex', gap:6, flexWrap:'wrap' }}>
          <span style={{ background:'rgba(0,0,0,.5)', color:'#fff', fontSize:11, fontWeight:600, padding:'2px 8px', borderRadius:20, display:'flex', alignItems:'center', gap:3 }}>
            <Ico.clock /> {t.duration}
          </span>
          <span style={{ background:'rgba(0,0,0,.5)', color:'#fff', fontSize:11, fontWeight:600, padding:'2px 8px', borderRadius:20, display:'flex', alignItems:'center', gap:3 }}>
            <Ico.car /> {t.transport}
          </span>
        </div>
      </div>
      <div style={{ padding:'14px 16px' }}>
        <RatingRow score={t.rating} count={t.reviews} size={11} />
        <div style={{ fontWeight:700, fontSize:15, marginTop:6, marginBottom:8, lineHeight:1.35, display:'-webkit-box', WebkitLineClamp:2, WebkitBoxOrient:'vertical', overflow:'hidden' }}>{t.name}</div>
        <div style={{ display:'flex', alignItems:'center', gap:5, color:'var(--text3)', fontSize:12, marginBottom:4 }}>
          <Ico.pin /> <span>{t.departure}</span>
          <span style={{ margin:'0 3px', color:'var(--border)' }}>→</span>
          <span>{t.dest.split(',')[0]}</span>
        </div>
        <div style={{ display:'flex', alignItems:'center', justifyContent:'space-between', marginTop:12 }}>
          <div>
            {t.oldPrice && <div style={{ fontSize:12, color:'var(--text3)', textDecoration:'line-through' }}>{formatPrice(t.oldPrice)}</div>}
            <div style={{ fontSize:16, fontWeight:800, color:'var(--accent)' }}>{formatPrice(t.price)}<span style={{ fontSize:12, fontWeight:500, color:'var(--text2)' }}>/người</span></div>
          </div>
          <button className="btn btn-primary btn-sm" onClick={e => { e.stopPropagation(); }}>Đặt ngay</button>
        </div>
      </div>
    </div>
  );
}

// ── Hotel Card ──
function HotelCard({ h, onClick }) {
  return (
    <div className="card" style={{ cursor:'pointer' }} onClick={onClick}>
      <div style={{ position:'relative' }}>
        <Img bg={h.bg} stripe={h.stripe} label={h.label} style={{ height:156, width:'100%' }} />
        <div style={{ position:'absolute', top:10, left:10 }}>
          <span style={{ background:'rgba(0,0,0,.55)', color:'var(--gold)', fontSize:12, fontWeight:700, padding:'3px 9px', borderRadius:20 }}>{'★'.repeat(h.stars)}</span>
        </div>
      </div>
      <div style={{ padding:'14px 16px' }}>
        <div style={{ fontWeight:700, fontSize:14, lineHeight:1.35, marginBottom:4, display:'-webkit-box', WebkitLineClamp:2, WebkitBoxOrient:'vertical', overflow:'hidden' }}>{h.name}</div>
        <div style={{ fontSize:12, color:'var(--text3)', marginBottom:8, display:'flex', gap:4, alignItems:'center' }}>
          <Ico.pin /> {h.loc} · {h.type}
        </div>
        <div style={{ display:'flex', gap:5, flexWrap:'wrap', marginBottom:10 }}>
          {h.amenities.slice(0,3).map(a => (
            <span key={a} style={{ fontSize:11, padding:'2px 7px', borderRadius:20, background:'var(--border-lt)', color:'var(--text2)', fontWeight:500 }}>{a}</span>
          ))}
          {h.amenities.length > 3 && <span style={{ fontSize:11, padding:'2px 7px', borderRadius:20, background:'var(--border-lt)', color:'var(--text3)' }}>+{h.amenities.length-3}</span>}
        </div>
        <div style={{ display:'flex', alignItems:'center', justifyContent:'space-between' }}>
          <div>
            <RatingRow score={h.rating} count={h.reviews} size={11} />
            <div style={{ fontSize:15, fontWeight:800, color:'var(--primary)', marginTop:3 }}>{formatPrice(h.price)}<span style={{ fontSize:11, fontWeight:500, color:'var(--text3)' }}>/đêm</span></div>
          </div>
          <button className="btn btn-outline btn-sm">Xem phòng</button>
        </div>
      </div>
    </div>
  );
}

// ── Combo Card ──
function ComboCard({ c, navigate }) {
  return (
    <div className="card" style={{ cursor:'pointer', display:'flex', overflow:'hidden' }} onClick={() => navigate('combo-detail',{combo:c})}>
      <Img bg={c.bg} stripe={c.stripe} label={c.label} style={{ width:160, flexShrink:0 }} />
      <div style={{ padding:'18px 20px', flex:1 }}>
        <div style={{ display:'flex', gap:8, marginBottom:10, flexWrap:'wrap' }}>
          <span className="badge badge-danger" style={{ fontSize:12, padding:'4px 10px' }}>Tiết kiệm {c.disc}%</span>
          <span className="badge badge-neutral">{c.duration}</span>
        </div>
        <div style={{ fontWeight:700, fontSize:15, lineHeight:1.35, marginBottom:8 }}>{c.name}</div>
        <div style={{ display:'flex', gap:5, marginBottom:10, flexWrap:'wrap' }}>
          {c.includes.map(inc => (
            <div key={inc} style={{ fontSize:12, color:'var(--text2)', display:'flex', gap:4, alignItems:'center' }}>
              <span style={{ color:'var(--success)' }}><Ico.check /></span> {inc}
            </div>
          ))}
        </div>
        <div style={{ fontSize:12, color:'var(--text3)', marginBottom:12, display:'flex', gap:4, alignItems:'center' }}>
          <Ico.plane /> Từ: {c.from}
        </div>
        <div style={{ display:'flex', alignItems:'center', justifyContent:'space-between' }}>
          <div>
            <div style={{ fontSize:12, color:'var(--text3)', textDecoration:'line-through' }}>{formatPrice(c.oldPrice)}</div>
            <div style={{ fontSize:18, fontWeight:800, color:'var(--accent)' }}>{formatPrice(c.price)}</div>
          </div>
          <button className="btn btn-accent btn-sm">Xem combo</button>
        </div>
      </div>
    </div>
  );
}

// ── Flash Sale Banner ──
function FlashSaleBanner({ navigate }) {
  const [h, setH] = useState(4); const [m, setM] = useState(23); const [s, setS] = useState(47);
  useEffect(() => {
    const t = setInterval(() => {
      setS(p => { if(p > 0) return p-1; setM(pm => { if(pm > 0) return pm-1; setH(ph => ph-1); return 59; }); return 59; });
    }, 1000);
    return () => clearInterval(t);
  }, []);
  return (
    <section style={{ background:'linear-gradient(135deg, var(--accent) 0%, oklch(0.55 0.22 20) 100%)', padding:'36px 0' }}>
      <div className="container" style={{ display:'flex', alignItems:'center', justifyContent:'space-between', gap:32, flexWrap:'wrap' }}>
        <div style={{ display:'flex', alignItems:'center', gap:24 }}>
          <div>
            <div style={{ fontSize:11, fontWeight:700, color:'rgba(255,255,255,.7)', letterSpacing:2, textTransform:'uppercase', marginBottom:4 }}>Flash Sale</div>
            <div style={{ fontSize:26, fontWeight:900, color:'#fff', letterSpacing:'-.02em' }}>Giảm đến 40%</div>
            <div style={{ fontSize:14, color:'rgba(255,255,255,.8)', marginTop:4 }}>Ưu đãi giới hạn — còn lại:</div>
          </div>
          <div style={{ display:'flex', gap:8, alignItems:'center' }}>
            {[h, m, s].map((v, i) => (
              <React.Fragment key={i}>
                <div style={{ background:'rgba(0,0,0,.25)', borderRadius:10, padding:'10px 14px', textAlign:'center', minWidth:60 }}>
                  <div style={{ fontSize:26, fontWeight:900, color:'#fff', lineHeight:1, fontVariantNumeric:'tabular-nums' }}>{String(v).padStart(2,'0')}</div>
                  <div style={{ fontSize:10, color:'rgba(255,255,255,.7)', marginTop:2 }}>{['GIỜ','PHÚT','GIÂY'][i]}</div>
                </div>
                {i < 2 && <span style={{ fontSize:22, fontWeight:900, color:'rgba(255,255,255,.7)', marginTop:-8 }}>:</span>}
              </React.Fragment>
            ))}
          </div>
        </div>
        <div style={{ display:'flex', gap:12, flexWrap:'wrap' }}>
          {TOURS.filter(t => t.oldPrice).slice(0,3).map(t => (
            <div key={t.id} style={{ background:'rgba(255,255,255,.15)', backdropFilter:'blur(8px)', borderRadius:'var(--r-md)', padding:'10px 16px', cursor:'pointer', transition:'background .15s', minWidth:170 }}>
              <div style={{ fontSize:12, fontWeight:700, color:'rgba(255,255,255,.8)', marginBottom:3 }}>{t.name.substring(0,28)}...</div>
              <div style={{ fontSize:13, color:'rgba(255,255,255,.6)', textDecoration:'line-through' }}>{formatPrice(t.oldPrice)}</div>
              <div style={{ fontSize:18, fontWeight:900, color:'#fff' }}>{formatPrice(t.price)}</div>
            </div>
          ))}
        </div>
        <button className="btn" style={{ background:'#fff', color:'var(--accent)', fontWeight:800, padding:'12px 28px', flexShrink:0 }} onClick={() => navigate('search',{type:'tour'})}>Mua ngay →</button>
      </div>
    </section>
  );
}

// ── Blog Card ──
function BlogCard({ b }) {
  return (
    <div className="card" style={{ cursor:'pointer', overflow:'hidden' }}>
      <Img bg={b.bg} stripe={b.stripe} label={b.label} style={{ height:180, width:'100%' }} />
      <div style={{ padding:'18px 20px' }}>
        <div style={{ display:'flex', gap:8, marginBottom:10 }}>
          <span className="badge badge-primary">{b.cat}</span>
          <span style={{ fontSize:12, color:'var(--text3)', display:'flex', alignItems:'center', gap:3 }}><Ico.clock /> {b.time}</span>
        </div>
        <div style={{ fontWeight:700, fontSize:16, lineHeight:1.4, marginBottom:8 }}>{b.title}</div>
        <p style={{ fontSize:13, color:'var(--text2)', lineHeight:1.6, display:'-webkit-box', WebkitLineClamp:3, WebkitBoxOrient:'vertical', overflow:'hidden' }}>{b.excerpt}</p>
        <div style={{ marginTop:14, display:'flex', justifyContent:'space-between', alignItems:'center' }}>
          <span style={{ fontSize:12, color:'var(--text3)' }}>{b.date}</span>
          <button className="btn btn-ghost btn-sm" style={{ gap:3, padding:'4px 10px' }}>Đọc thêm <Ico.chevR /></button>
        </div>
      </div>
    </div>
  );
}

// ── Footer ──
function Footer() {
  const cols = [
    { title:'Về TravelViet', links:['Giới thiệu','Điều khoản sử dụng','Chính sách bảo mật','Chính sách hoàn hủy','Tuyển dụng'] },
    { title:'Khám phá', links:['Tour nội địa','Khách sạn','Combo du lịch','Blog du lịch','Khuyến mãi'] },
    { title:'Hỗ trợ', links:['Hướng dẫn đặt tour','Hướng dẫn đặt KS','Câu hỏi thường gặp','Liên hệ hỗ trợ','Trung tâm trợ giúp'] },
    { title:'Đối tác', links:['Trở thành đại lý','Chương trình Affiliate','Corporate Booking','API tích hợp','Hợp tác truyền thông'] },
  ];
  return (
    <footer style={{ background:'var(--text)', color:'rgba(255,255,255,.75)', padding:'48px 0 24px' }}>
      <div className="container">
        <div style={{ display:'grid', gridTemplateColumns:'2fr 1fr 1fr 1fr', gap:40, marginBottom:40 }}>
          <div>
            <div style={{ fontSize:22, fontWeight:900, color:'#fff', letterSpacing:'-.02em', marginBottom:12 }}>Travel<span style={{ color:'var(--accent)' }}>Viet</span></div>
            <p style={{ fontSize:13, lineHeight:1.7, color:'rgba(255,255,255,.6)', maxWidth:280 }}>Nền tảng đặt tour & khách sạn nội địa hàng đầu Việt Nam. Kết nối hàng nghìn đại lý uy tín với khách hàng trên cả nước.</p>
            <div style={{ marginTop:16, display:'flex', gap:10 }}>
              {['Facebook','Zalo','Instagram','YouTube'].map(s => (
                <div key={s} style={{ width:34, height:34, borderRadius:8, background:'rgba(255,255,255,.1)', display:'flex', alignItems:'center', justifyContent:'center', fontSize:12, fontWeight:700, cursor:'pointer', color:'rgba(255,255,255,.7)' }}>{s[0]}</div>
              ))}
            </div>
            <div style={{ marginTop:16, fontSize:13 }}>
              <div style={{ color:'rgba(255,255,255,.5)', marginBottom:4 }}>Hotline hỗ trợ</div>
              <div style={{ color:'#fff', fontWeight:700, fontSize:16 }}>1800 6006</div>
            </div>
          </div>
          {cols.slice(1).map(col => (
            <div key={col.title}>
              <div style={{ fontSize:13, fontWeight:700, color:'#fff', letterSpacing:.5, marginBottom:14, textTransform:'uppercase' }}>{col.title}</div>
              <div style={{ display:'flex', flexDirection:'column', gap:8 }}>
                {col.links.map(l => <a key={l} href="#" style={{ fontSize:13, color:'rgba(255,255,255,.6)', transition:'color .15s' }} onMouseEnter={e=>e.target.style.color='#fff'} onMouseLeave={e=>e.target.style.color='rgba(255,255,255,.6)'}>{l}</a>)}
              </div>
            </div>
          ))}
        </div>
        <div style={{ borderTop:'1px solid rgba(255,255,255,.1)', paddingTop:20, display:'flex', justifyContent:'space-between', alignItems:'center', fontSize:12, color:'rgba(255,255,255,.4)' }}>
          <span>© 2026 TravelViet. Bảo lưu mọi quyền.</span>
          <span>Giấy phép kinh doanh: 0123456789 — Bộ KHĐT</span>
        </div>
      </div>
    </footer>
  );
}

Object.assign(window, { HomePage, TourCard, HotelCard, Footer });
