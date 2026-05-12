// DealsPage — Promotion list + Deal detail
const DEALS_DATA = [
  { id:1, code:'SUMMER30', title:'Flash Sale Hè 2026', type:'flash', disc:'30%', discType:'percent', discVal:30, minOrder:2000000, maxDisc:500000, expiry:'30/06/2026', applicable:'Tất cả tour & khách sạn', used:342, total:500, bg:'#c2410c', stripe:'#9a3412', tags:['Hot','Flash Sale'] },
  { id:2, code:'HALONG200K', title:'Ưu đãi Hạ Long Bay', type:'location', disc:'200.000đ', discType:'fixed', discVal:200000, minOrder:1500000, maxDisc:200000, expiry:'15/06/2026', applicable:'Tour Hạ Long', used:89, total:200, bg:'#0c4a6e', stripe:'#083b5a', tags:['Địa điểm'] },
  { id:3, code:'EARLYBIRD15', title:'Đặt sớm — Nghỉ thỏa thích', type:'earlybird', disc:'15%', discType:'percent', discVal:15, minOrder:3000000, maxDisc:600000, expiry:'31/08/2026', applicable:'Đặt trước 30 ngày', used:156, total:1000, bg:'#2d6a4f', stripe:'#1b4332', tags:['Early Bird'] },
  { id:4, code:'NEWUSER50K', title:'Chào mừng thành viên mới', type:'newuser', disc:'50.000đ', discType:'fixed', discVal:50000, minOrder:500000, maxDisc:50000, expiry:'31/12/2026', applicable:'Lần đặt đầu tiên', used:2341, total:9999, bg:'#5b21b6', stripe:'#3d1580', tags:['Thành viên mới'] },
  { id:5, code:'GOLD300K', title:'Đặc quyền hạng Vàng', type:'vip', disc:'300.000đ', discType:'fixed', discVal:300000, minOrder:5000000, maxDisc:300000, expiry:'31/12/2026', applicable:'Hạng Vàng trở lên', used:45, total:500, bg:'#92400e', stripe:'#6b2f09', tags:['VIP'] },
  { id:6, code:'PHUQUOC20', title:'Phú Quốc Sale lớn', type:'location', disc:'20%', discType:'percent', discVal:20, minOrder:4000000, maxDisc:800000, expiry:'30/06/2026', applicable:'Tour & KS Phú Quốc', used:167, total:300, bg:'#0369a1', stripe:'#024d79', tags:['Hot','Địa điểm'] },
  { id:7, code:'LASTMIN10', title:'Last Minute Deal', type:'lastminute', disc:'10%', discType:'percent', discVal:10, minOrder:1000000, maxDisc:400000, expiry:'05/06/2026', applicable:'Đặt trong 3 ngày tới', used:78, total:200, bg:'#0e7490', stripe:'#0a5570', tags:['Last Minute'] },
  { id:8, code:'DALAT15', title:'Đà Lạt Mùa Hoa', type:'location', disc:'15%', discType:'percent', discVal:15, minOrder:1500000, maxDisc:300000, expiry:'30/06/2026', applicable:'Tour & KS Đà Lạt', used:123, total:400, bg:'#5b21b6', stripe:'#3d1580', tags:['Địa điểm'] },
];

function DealsPage({ navigate }) {
  const [activeFilter, setActiveFilter] = useState('all');
  const [selDeal, setSelDeal] = useState(null);
  const [copied, setCopied] = useState({});
  const [email, setEmail] = useState('');
  const [subscribed, setSubscribed] = useState(false);

  function copyCode(code, e) {
    e?.stopPropagation();
    navigator.clipboard?.writeText(code).catch(() => {});
    setCopied(p => ({ ...p, [code]:true }));
    setTimeout(() => setCopied(p => ({ ...p, [code]:false })), 2500);
  }

  const filters = [
    { id:'all', label:'Tất cả', count: DEALS_DATA.length },
    { id:'flash', label:'Flash Sale' },
    { id:'earlybird', label:'Early Bird' },
    { id:'location', label:'Theo điểm đến' },
    { id:'vip', label:'Thành viên VIP' },
    { id:'lastminute', label:'Last Minute' },
    { id:'newuser', label:'Thành viên mới' },
  ];

  const filtered = activeFilter === 'all' ? DEALS_DATA : DEALS_DATA.filter(d => d.type === activeFilter);

  if (selDeal) return <DealDetail deal={selDeal} onBack={() => setSelDeal(null)} navigate={navigate} copied={copied} copyCode={copyCode} />;

  return (
    <div className="page-anim" style={{ paddingTop:68 }}>
      {/* Hero banner */}
      <div style={{ background:'linear-gradient(135deg, var(--primary) 0%, oklch(0.27 0.16 240) 100%)', padding:'52px 0 64px', position:'relative', overflow:'hidden' }}>
        <div style={{ position:'absolute', top:-60, right:-60, width:360, height:360, borderRadius:'50%', background:'rgba(255,255,255,.04)' }} />
        <div style={{ position:'absolute', bottom:-80, left:80, width:240, height:240, borderRadius:'50%', background:'rgba(255,255,255,.03)' }} />
        <div style={{ position:'absolute', top:20, right:200, width:140, height:140, borderRadius:'50%', background:'rgba(255,255,255,.04)' }} />
        <div className="container" style={{ position:'relative', zIndex:1 }}>
          <div style={{ display:'flex', gap:8, marginBottom:14 }}>
            <span style={{ background:'var(--accent)', color:'#fff', padding:'4px 14px', borderRadius:'var(--r-full)', fontSize:12, fontWeight:700, letterSpacing:.5, whiteSpace:'nowrap' }}>KHUYẾN MÃI & ƯU ĐÃI</span>
          </div>
          <h1 style={{ fontSize:44, fontWeight:900, color:'#fff', letterSpacing:'-.02em', lineHeight:1.15, marginBottom:12, maxWidth:560, textWrap:'balance' }}>
            Ưu đãi hấp dẫn<br /><span style={{ color:'rgba(255,255,255,.7)' }}>cập nhật mỗi ngày</span>
          </h1>
          <p style={{ fontSize:16, color:'rgba(255,255,255,.78)', maxWidth:480, marginBottom:28, lineHeight:1.6 }}>
            Hàng chục voucher và deal du lịch độc quyền. Tiết kiệm đến 30% cho tour, khách sạn và combo trên toàn Việt Nam.
          </p>
          <div style={{ display:'flex', gap:32 }}>
            {[['8','Deal đang chạy'],['30%','Giảm tối đa'],['500','Lượt dùng hôm nay'],['4.9★','Hài lòng khách hàng']].map(([v,l]) => (
              <div key={l} style={{ textAlign:'center' }}>
                <div style={{ fontSize:26, fontWeight:900, color:'#fff', lineHeight:1 }}>{v}</div>
                <div style={{ fontSize:12, color:'rgba(255,255,255,.65)', marginTop:4 }}>{l}</div>
              </div>
            ))}
          </div>
        </div>
      </div>

      <div className="container" style={{ marginTop:'-24px', paddingBottom:56 }}>
        {/* Filter bar */}
        <div style={{ background:'#fff', borderRadius:'var(--r-lg)', padding:'16px 20px', boxShadow:'var(--sh-md)', marginBottom:28, display:'flex', gap:8, flexWrap:'wrap', alignItems:'center' }}>
          <span style={{ fontSize:13, fontWeight:700, color:'var(--text2)', marginRight:4 }}>Lọc theo:</span>
          {filters.map(f => (
            <button key={f.id} onClick={() => setActiveFilter(f.id)}
              className={'btn btn-sm ' + (activeFilter===f.id ? 'btn-primary' : 'btn-ghost')}
              style={{ borderRadius:'var(--r-full)', padding:'6px 14px' }}>
              {f.label}{f.count ? ` (${f.count})` : ''}
            </button>
          ))}
          <div style={{ marginLeft:'auto', fontSize:13, color:'var(--text3)' }}>{filtered.length} ưu đãi</div>
        </div>

        {/* Flash sale highlight */}
        {activeFilter === 'all' && (
          <div style={{ background:'linear-gradient(135deg, var(--accent) 0%, oklch(0.52 0.22 20) 100%)', borderRadius:'var(--r-xl)', padding:'24px 28px', marginBottom:28, display:'flex', alignItems:'center', justifyContent:'space-between', gap:20, flexWrap:'wrap' }}>
            <div>
              <div style={{ fontSize:12, fontWeight:700, color:'rgba(255,255,255,.7)', letterSpacing:1.5, textTransform:'uppercase', marginBottom:6 }}>🔥 Flash Sale đang diễn ra</div>
              <div style={{ fontSize:24, fontWeight:900, color:'#fff', marginBottom:4 }}>Giảm 30% toàn bộ tour & KS</div>
              <div style={{ fontSize:14, color:'rgba(255,255,255,.8)' }}>Chỉ còn {DEALS_DATA[0].total - DEALS_DATA[0].used} lượt · HSD 30/06/2026</div>
            </div>
            <div style={{ display:'flex', gap:12, alignItems:'center' }}>
              <div style={{ background:'rgba(0,0,0,.2)', borderRadius:'var(--r-md)', padding:'10px 20px', textAlign:'center', fontFamily:'monospace' }}>
                <div style={{ fontSize:28, fontWeight:900, color:'#fff', lineHeight:1 }}>{DEALS_DATA[0].code}</div>
                <div style={{ fontSize:11, color:'rgba(255,255,255,.7)', marginTop:2 }}>MÃ VOUCHER</div>
              </div>
              <button onClick={(e) => copyCode(DEALS_DATA[0].code, e)}
                className="btn btn-lg" style={{ background:'#fff', color:'var(--accent)', fontWeight:800, flexShrink:0 }}>
                {copied[DEALS_DATA[0].code] ? '✓ Đã copy!' : 'Copy mã ngay'}
              </button>
            </div>
          </div>
        )}

        {/* Deals grid */}
        <div style={{ display:'grid', gridTemplateColumns:'repeat(4,1fr)', gap:18, marginBottom:40 }}>
          {filtered.map(deal => (
            <DealCard key={deal.id} deal={deal} copied={copied} copyCode={copyCode} onDetail={() => setSelDeal(deal)} />
          ))}
        </div>

        {/* Subscribe section */}
        <div style={{ background:'var(--primary-xlt)', borderRadius:'var(--r-xl)', padding:'36px', textAlign:'center', border:'1px solid var(--primary-lt)' }}>
          <div style={{ fontSize:22, fontWeight:800, marginBottom:8 }}>Nhận thông báo deal mới</div>
          <p style={{ fontSize:14, color:'var(--text2)', marginBottom:20, maxWidth:440, margin:'0 auto 20px' }}>
            Đăng ký để nhận ngay thông báo khi có ưu đãi mới cho điểm đến yêu thích của bạn.
          </p>
          {subscribed ? (
            <div style={{ fontSize:15, fontWeight:700, color:'var(--success)', display:'flex', gap:8, alignItems:'center', justifyContent:'center' }}>
              <Ico.check /> Đăng ký thành công! Chúng tôi sẽ thông báo sớm.
            </div>
          ) : (
            <div style={{ display:'flex', gap:10, justifyContent:'center', maxWidth:400, margin:'0 auto' }}>
              <input className="input" placeholder="Email của bạn..." value={email} onChange={e => setEmail(e.target.value)} style={{ flex:1 }} />
              <button onClick={() => email && setSubscribed(true)} className="btn btn-primary" style={{ flexShrink:0 }}>Đăng ký</button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

function DealCard({ deal, copied, copyCode, onDetail }) {
  const pct = Math.min(100, Math.round(deal.used / deal.total * 100));
  const left = deal.total - deal.used;
  const urgent = left < 50;

  return (
    <div style={{ background:'#fff', borderRadius:'var(--r-lg)', boxShadow:'var(--sh-sm)', overflow:'hidden', cursor:'pointer', transition:'transform .2s, box-shadow .2s' }}
      onMouseEnter={e => { e.currentTarget.style.transform='translateY(-4px)'; e.currentTarget.style.boxShadow='var(--sh-md)'; }}
      onMouseLeave={e => { e.currentTarget.style.transform='none'; e.currentTarget.style.boxShadow='var(--sh-sm)'; }}
      onClick={onDetail}>

      {/* Color banner with pattern */}
      <div style={{ background:deal.bg, height:76, position:'relative', overflow:'hidden' }}>
        <svg width="100%" height="100%" style={{ position:'absolute', inset:0 }}>
          <defs>
            <pattern id={'dp'+deal.id} width="24" height="24" patternUnits="userSpaceOnUse" patternTransform="rotate(42)">
              <rect width="12" height="24" fill={deal.stripe} opacity=".28"/>
            </pattern>
          </defs>
          <rect width="100%" height="100%" fill={'url(#dp'+deal.id+')'}/>
        </svg>
        <div style={{ position:'absolute', inset:0, display:'flex', alignItems:'center', justifyContent:'space-between', padding:'0 16px' }}>
          <div style={{ fontSize:30, fontWeight:900, color:'#fff', letterSpacing:-.5 }}>{deal.disc}</div>
          <div style={{ display:'flex', gap:4, flexWrap:'wrap', justifyContent:'flex-end', maxWidth:90 }}>
            {deal.tags.map(t => <span key={t} style={{ fontSize:10, padding:'1px 6px', borderRadius:20, background:'rgba(255,255,255,.2)', color:'#fff', fontWeight:700, whiteSpace:'nowrap' }}>{t}</span>)}
          </div>
        </div>
      </div>

      {/* Content */}
      <div style={{ padding:'14px 16px 16px' }}>
        <div style={{ fontWeight:700, fontSize:14, marginBottom:4, lineHeight:1.3 }}>{deal.title}</div>
        <div style={{ fontSize:12, color:'var(--text3)', marginBottom:10 }}>{deal.applicable}</div>

        {/* Usage progress */}
        <div style={{ marginBottom:12 }}>
          <div style={{ display:'flex', justifyContent:'space-between', fontSize:11, color:'var(--text3)', marginBottom:4 }}>
            <span>Đã dùng {deal.used.toLocaleString()}</span>
            <span style={{ color:urgent ? 'var(--danger)' : 'var(--text3)', fontWeight: urgent ? 700 : 400 }}>
              {urgent ? '⚡ ' : ''}Còn {left.toLocaleString()}
            </span>
          </div>
          <div style={{ height:6, background:'var(--border-lt)', borderRadius:3, overflow:'hidden' }}>
            <div style={{ height:'100%', width:pct+'%', background: pct>85 ? 'var(--danger)' : pct>60 ? 'var(--warn)' : 'var(--accent)', borderRadius:3, transition:'width .5s' }} />
          </div>
        </div>

        {/* Code + copy */}
        <div style={{ display:'flex', gap:6, marginBottom:10 }}>
          <div style={{ flex:1, padding:'8px 10px', background:'var(--bg)', borderRadius:'var(--r-sm)', border:'1.5px dashed var(--border)', fontSize:13, fontWeight:800, color:'var(--primary)', fontFamily:'monospace', letterSpacing:.8, overflow:'hidden', textOverflow:'ellipsis', whiteSpace:'nowrap' }}>{deal.code}</div>
          <button onClick={(e) => copyCode(deal.code, e)}
            style={{ padding:'8px 12px', borderRadius:'var(--r-sm)', border:'none', cursor:'pointer', fontWeight:700, fontSize:12, transition:'all .15s', flexShrink:0,
              background: copied[deal.code] ? 'var(--success)' : 'var(--primary)',
              color:'#fff' }}>
            {copied[deal.code] ? '✓' : 'Copy'}
          </button>
        </div>

        <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center' }}>
          <span style={{ fontSize:11, color:'var(--text3)' }}>HSD: {deal.expiry}</span>
          <button onClick={onDetail} style={{ background:'none', border:'none', cursor:'pointer', fontSize:12, color:'var(--primary)', fontWeight:600, display:'flex', gap:3, alignItems:'center', padding:'0' }}>
            Chi tiết <Ico.chevR />
          </button>
        </div>
      </div>
    </div>
  );
}

function DealDetail({ deal, onBack, navigate, copied, copyCode }) {
  const left = deal.total - deal.used;
  return (
    <div className="page-anim" style={{ paddingTop:68 }}>
      {/* Breadcrumb */}
      <div style={{ borderBottom:'1px solid var(--border-lt)', background:'var(--surface)' }}>
        <div className="container" style={{ padding:'12px 32px', display:'flex', gap:6, fontSize:13, color:'var(--text3)', alignItems:'center' }}>
          <button onClick={() => navigate('home')} style={{ background:'none',border:'none',cursor:'pointer',color:'var(--primary)',fontSize:13,padding:0,fontWeight:500 }}>Trang chủ</button>
          <Ico.chevR />
          <button onClick={onBack} style={{ background:'none',border:'none',cursor:'pointer',color:'var(--primary)',fontSize:13,padding:0,fontWeight:500 }}>Khuyến mãi</button>
          <Ico.chevR />
          <span style={{ color:'var(--text)', fontWeight:500 }}>{deal.title}</span>
        </div>
      </div>

      <div className="container" style={{ paddingTop:32, paddingBottom:48 }}>
        <div style={{ display:'grid', gridTemplateColumns:'1fr 360px', gap:28, alignItems:'flex-start' }}>
          {/* Left content */}
          <div style={{ display:'flex', flexDirection:'column', gap:20 }}>
            {/* Main hero card */}
            <div style={{ background:deal.bg, borderRadius:'var(--r-xl)', padding:'36px', position:'relative', overflow:'hidden' }}>
              <svg width="100%" height="100%" style={{ position:'absolute', inset:0 }}>
                <defs><pattern id={'dd'+deal.id} width="28" height="28" patternUnits="userSpaceOnUse" patternTransform="rotate(42)"><rect width="14" height="28" fill={deal.stripe} opacity=".28"/></pattern></defs>
                <rect width="100%" height="100%" fill={'url(#dd'+deal.id+')'}/>
              </svg>
              <div style={{ position:'relative' }}>
                <div style={{ display:'flex', gap:6, marginBottom:14, flexWrap:'wrap' }}>
                  {deal.tags.map(t => <span key={t} style={{ fontSize:11, padding:'3px 10px', borderRadius:20, background:'rgba(255,255,255,.2)', color:'#fff', fontWeight:700 }}>{t}</span>)}
                </div>
                <div style={{ fontSize:56, fontWeight:900, color:'#fff', lineHeight:1, marginBottom:8 }}>{deal.disc}</div>
                <div style={{ fontSize:22, fontWeight:700, color:'rgba(255,255,255,.9)', marginBottom:6 }}>{deal.title}</div>
                <div style={{ fontSize:14, color:'rgba(255,255,255,.75)' }}>Áp dụng cho: {deal.applicable}</div>
              </div>
            </div>

            {/* Detail grid */}
            <div style={{ background:'#fff', borderRadius:'var(--r-lg)', padding:'24px', border:'1px solid var(--border-lt)' }}>
              <h2 style={{ fontWeight:800, fontSize:18, marginBottom:18 }}>Chi tiết ưu đãi</h2>
              <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:12 }}>
                {[
                  ['Mức giảm', deal.disc],
                  ['Đơn tối thiểu', formatPrice(deal.minOrder)],
                  ['Giảm tối đa', formatPrice(deal.maxDisc)],
                  ['Hạn sử dụng', deal.expiry],
                  ['Tổng lượt', deal.total.toLocaleString() + ' voucher'],
                  ['Đã sử dụng', deal.used.toLocaleString() + ' lượt'],
                ].map(([k,v]) => (
                  <div key={k} style={{ padding:'14px 16px', background:'var(--bg)', borderRadius:'var(--r-md)' }}>
                    <div style={{ fontSize:11, color:'var(--text3)', fontWeight:700, marginBottom:4, textTransform:'uppercase', letterSpacing:.4 }}>{k}</div>
                    <div style={{ fontSize:15, fontWeight:700, color:'var(--text)' }}>{v}</div>
                  </div>
                ))}
              </div>
              {/* Usage bar */}
              <div style={{ marginTop:16 }}>
                <div style={{ display:'flex', justifyContent:'space-between', fontSize:12, color:'var(--text3)', marginBottom:6 }}>
                  <span>Đã dùng {deal.used.toLocaleString()} / {deal.total.toLocaleString()}</span>
                  <span style={{ color: left < 50 ? 'var(--danger)' : 'var(--text2)', fontWeight:700 }}>Còn {left.toLocaleString()} lượt</span>
                </div>
                <div style={{ height:10, background:'var(--border-lt)', borderRadius:5, overflow:'hidden' }}>
                  <div style={{ height:'100%', width:Math.min(100,deal.used/deal.total*100)+'%', background: deal.used/deal.total > 0.85 ? 'var(--danger)' : 'var(--accent)', borderRadius:5 }} />
                </div>
              </div>
            </div>

            {/* Terms */}
            <div style={{ background:'#fff', borderRadius:'var(--r-lg)', padding:'24px', border:'1px solid var(--border-lt)' }}>
              <h2 style={{ fontWeight:800, fontSize:18, marginBottom:16 }}>Điều kiện & Điều khoản</h2>
              <ul style={{ listStyle:'none', display:'flex', flexDirection:'column', gap:10 }}>
                {[
                  'Mỗi tài khoản chỉ sử dụng mã này 1 lần duy nhất',
                  'Không áp dụng cộng dồn với voucher khuyến mãi khác',
                  `Đơn hàng tối thiểu ${formatPrice(deal.minOrder)} (trước khi giảm)`,
                  `Mức giảm tối đa ${formatPrice(deal.maxDisc)} (áp dụng cho loại giảm %)`,
                  `Chỉ áp dụng cho: ${deal.applicable}`,
                  `Voucher hết hiệu lực lúc 23:59 ngày ${deal.expiry}`,
                  'TravelViet có quyền thu hồi nếu phát hiện gian lận hoặc lạm dụng',
                  'Không áp dụng cho các booking đã thanh toán trước đó',
                ].map((t,i) => (
                  <li key={i} style={{ display:'flex', gap:10, fontSize:14, color:'var(--text2)', alignItems:'flex-start', lineHeight:1.6 }}>
                    <span style={{ color:'var(--primary)', fontWeight:700, flexShrink:0, marginTop:2 }}>•</span>{t}
                  </li>
                ))}
              </ul>
            </div>
          </div>

          {/* Right sticky */}
          <div style={{ position:'sticky', top:84 }}>
            <div className="card-flat" style={{ borderRadius:'var(--r-xl)', padding:'28px', border:'2px solid var(--border)' }}>
              <div style={{ fontWeight:800, fontSize:17, marginBottom:16 }}>Sử dụng ngay</div>

              {/* Voucher code display */}
              <div style={{ padding:'18px', background:'var(--bg)', borderRadius:'var(--r-lg)', textAlign:'center', marginBottom:14, border:'2px dashed var(--border)' }}>
                <div style={{ fontSize:11, color:'var(--text3)', fontWeight:700, marginBottom:8, letterSpacing:1, textTransform:'uppercase' }}>Mã voucher của bạn</div>
                <div style={{ fontSize:28, fontWeight:900, color:'var(--primary)', letterSpacing:4, fontFamily:'monospace', marginBottom:4 }}>{deal.code}</div>
                <div style={{ fontSize:12, color:'var(--text3)' }}>Giảm {deal.disc} · Tối đa {formatPrice(deal.maxDisc)}</div>
              </div>

              <button onClick={(e) => copyCode(deal.code, e)}
                className={'btn btn-lg ' + (copied[deal.code] ? 'btn-primary' : 'btn-accent')}
                style={{ width:'100%', justifyContent:'center', marginBottom:14, letterSpacing:.5 }}>
                {copied[deal.code] ? '✓ Đã sao chép mã voucher!' : 'Sao chép mã giảm giá'}
              </button>

              <hr style={{ marginBottom:14 }} />
              <div style={{ fontWeight:700, fontSize:14, marginBottom:12 }}>Áp dụng ngay cho</div>
              <div style={{ display:'flex', flexDirection:'column', gap:8 }}>
                {[
                  ['search',{type:'tour'},'✈ Tìm tour du lịch'],
                  ['search',{type:'hotel'},'🏨 Tìm khách sạn'],
                  ['search',{type:'combo'},'📦 Tìm combo tiết kiệm'],
                ].map(([page, params, label]) => (
                  <button key={label} onClick={() => navigate(page, params)}
                    className="btn btn-outline btn-sm" style={{ justifyContent:'flex-start', gap:8 }}>{label}</button>
                ))}
              </div>

              <div style={{ marginTop:16, padding:'12px 14px', background:'var(--danger-lt)', borderRadius:'var(--r-md)', fontSize:12, color:'var(--text2)', lineHeight:1.6 }}>
                <strong style={{ color:'var(--danger)' }}>⚡ Sắp hết:</strong> Chỉ còn <strong>{left}</strong> lượt · Hết hạn {deal.expiry}
              </div>
            </div>

            {/* Related deals */}
            <div style={{ marginTop:16 }}>
              <div style={{ fontWeight:700, fontSize:14, marginBottom:12, color:'var(--text)' }}>Ưu đãi tương tự</div>
              <div style={{ display:'flex', flexDirection:'column', gap:10 }}>
                {DEALS_DATA.filter(d => d.id !== deal.id).slice(0,3).map(d => (
                  <div key={d.id} onClick={() => { window.scrollTo({top:0,behavior:'smooth'}); setTimeout(()=>{ /* handled by parent */ },100); }}
                    style={{ background:'#fff', borderRadius:'var(--r-md)', padding:'12px 14px', border:'1px solid var(--border-lt)', cursor:'pointer', display:'flex', justifyContent:'space-between', alignItems:'center', transition:'all .15s' }}
                    onMouseEnter={e => e.currentTarget.style.borderColor='var(--primary)'}
                    onMouseLeave={e => e.currentTarget.style.borderColor='var(--border-lt)'}>
                    <div>
                      <div style={{ fontWeight:600, fontSize:13 }}>{d.title}</div>
                      <div style={{ fontSize:12, color:'var(--text3)' }}>{d.expiry}</div>
                    </div>
                    <span style={{ fontWeight:900, fontSize:16, color:'var(--accent)' }}>{d.disc}</span>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

Object.assign(window, { DealsPage });
