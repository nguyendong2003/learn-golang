// AccountPage — User dashboard: bookings, loyalty, profile, vouchers
function AccountPage({ navigate }) {
  const [activeMenu, setActiveMenu] = useState('bookings');
  const [filterStatus, setFilterStatus] = useState('all');

  const menus = [
    { id:'bookings', label:'Đặt chỗ của tôi', icon:'📋' },
    { id:'loyalty', label:'Điểm thưởng & Hạng', icon:'⭐' },
    { id:'vouchers', label:'Voucher của tôi', icon:'🎁' },
    { id:'reviews', label:'Đánh giá của tôi', icon:'💬' },
    { id:'affiliate', label:'Cộng tác viên', icon:'🔗' },
    { id:'profile', label:'Thông tin cá nhân', icon:'👤' },
    { id:'security', label:'Bảo mật tài khoản', icon:'🔒' },
  ];

  const statusMap = {
    confirmed: { label:'Đã xác nhận', cls:'st-confirmed' },
    paid: { label:'Đã thanh toán', cls:'st-paid' },
    completed: { label:'Hoàn tất', cls:'st-completed' },
    cancelled: { label:'Đã hủy', cls:'st-cancelled' },
    pending: { label:'Chờ xác nhận', cls:'st-pending' },
    inprogress: { label:'Đang sử dụng', cls:'st-inprogress' },
  };

  const filtered = filterStatus === 'all' ? BOOKINGS : BOOKINGS.filter(b => b.status === filterStatus);

  return (
    <div className="page-anim" style={{ paddingTop:68, background:'var(--bg)', minHeight:'100vh' }}>
      {/* ── Profile header ── */}
      <div style={{ background:'linear-gradient(135deg, var(--primary) 0%, oklch(0.30 0.14 235) 100%)', padding:'36px 0 60px' }}>
        <div className="container">
          <div style={{ display:'flex', gap:20, alignItems:'center' }}>
            <div style={{ width:72, height:72, borderRadius:'50%', background:'rgba(255,255,255,.2)', border:'3px solid rgba(255,255,255,.4)', display:'flex', alignItems:'center', justifyContent:'center', fontSize:28, fontWeight:900, color:'#fff', flexShrink:0 }}>NT</div>
            <div style={{ flex:1 }}>
              <div style={{ fontSize:22, fontWeight:900, color:'#fff', marginBottom:4 }}>{USER.name}</div>
              <div style={{ display:'flex', gap:10, flexWrap:'wrap', alignItems:'center' }}>
                <TierBadge tier={USER.tier} />
                <span style={{ fontSize:13, color:'rgba(255,255,255,.75)' }}>{USER.points.toLocaleString()} điểm thưởng</span>
                <span style={{ fontSize:13, color:'rgba(255,255,255,.5)' }}>·</span>
                <span style={{ fontSize:13, color:'rgba(255,255,255,.75)' }}>Thành viên từ {USER.joinDate}</span>
              </div>
            </div>
            <div style={{ display:'flex', gap:24, textAlign:'center' }}>
              {[
                { val: BOOKINGS.length, lbl:'Tổng booking' },
                { val: BOOKINGS.filter(b=>b.status==='completed').length, lbl:'Đã hoàn tất' },
                { val: formatPrice(USER.totalSpend), lbl:'Tổng chi tiêu' },
              ].map(({val, lbl}) => (
                <div key={lbl}>
                  <div style={{ fontSize:22, fontWeight:900, color:'#fff' }}>{val}</div>
                  <div style={{ fontSize:12, color:'rgba(255,255,255,.65)' }}>{lbl}</div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>

      <div className="container" style={{ marginTop:'-28px', paddingBottom:48 }}>
        <div style={{ display:'flex', gap:24, alignItems:'flex-start' }}>

          {/* ── Sidebar menu ── */}
          <aside style={{ width:220, flexShrink:0, position:'sticky', top:84 }}>
            <div style={{ background:'#fff', borderRadius:'var(--r-lg)', boxShadow:'var(--sh-sm)', overflow:'hidden', border:'1px solid var(--border-lt)' }}>
              {menus.map((m, i) => (
                <button key={m.id} onClick={() => setActiveMenu(m.id)}
                  style={{ width:'100%', display:'flex', alignItems:'center', gap:10, padding:'13px 18px', border:'none', cursor:'pointer', fontSize:14, fontWeight: activeMenu===m.id ? 700 : 500, transition:'all .15s', textAlign:'left', borderLeft: activeMenu===m.id ? '3px solid var(--primary)' : '3px solid transparent', background: activeMenu===m.id ? 'var(--primary-xlt)' : 'transparent', color: activeMenu===m.id ? 'var(--primary)' : 'var(--text2)', borderTop: i>0 ? '1px solid var(--border-lt)' : 'none' }}>
                  <span style={{ fontSize:16 }}>{m.icon}</span> {m.label}
                </button>
              ))}
            </div>
          </aside>

          {/* ── Main content ── */}
          <div style={{ flex:1, minWidth:0 }}>

            {/* BOOKINGS */}
            {activeMenu === 'bookings' && (
              <div>
                <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center', marginBottom:20, background:'#fff', padding:'18px 22px', borderRadius:'var(--r-lg)', border:'1px solid var(--border-lt)' }}>
                  <h2 style={{ fontWeight:800, fontSize:18 }}>Đặt chỗ của tôi</h2>
                  <div style={{ display:'flex', gap:6 }}>
                    {[['all','Tất cả'],['confirmed','Đã xác nhận'],['paid','Đã thanh toán'],['completed','Hoàn tất'],['cancelled','Đã hủy']].map(([id,lbl]) => (
                      <button key={id} onClick={() => setFilterStatus(id)}
                        className={'btn btn-sm' + (filterStatus===id ? ' btn-primary' : ' btn-ghost')} style={{ padding:'6px 12px' }}>{lbl}</button>
                    ))}
                  </div>
                </div>
                <div style={{ display:'flex', flexDirection:'column', gap:14 }}>
                  {filtered.length === 0 ? (
                    <div style={{ textAlign:'center', padding:'48px', background:'#fff', borderRadius:'var(--r-lg)', color:'var(--text3)' }}>
                      <div style={{ fontSize:40, marginBottom:10 }}>📭</div>
                      <div style={{ fontWeight:600 }}>Không có đặt chỗ nào</div>
                    </div>
                  ) : filtered.map(b => <BookingCard key={b.id} b={b} statusMap={statusMap} navigate={navigate} />)}
                </div>
              </div>
            )}

            {/* LOYALTY */}
            {activeMenu === 'loyalty' && <LoyaltyPanel />}

            {/* VOUCHERS */}
            {activeMenu === 'vouchers' && <VouchersPanel />}

            {/* REVIEWS */}
            {activeMenu === 'reviews' && <ReviewsPanel />}

            {/* AFFILIATE */}
            {activeMenu === 'affiliate' && <AffiliatePanel />}

            {/* PROFILE */}
            {activeMenu === 'profile' && <ProfilePanel />}

            {/* SECURITY */}
            {activeMenu === 'security' && <SecurityPanel />}
          </div>
        </div>
      </div>
    </div>
  );
}

function BookingCard({ b, statusMap, navigate }) {
  const st = statusMap[b.status] || { label: b.statusLabel, cls:'st-pending' };
  const typeIcon = { tour:'✈', hotel:'🏨', combo:'📦' };
  return (
    <div style={{ background:'#fff', borderRadius:'var(--r-lg)', border:'1px solid var(--border-lt)', overflow:'hidden' }}>
      <div style={{ display:'flex', gap:0 }}>
        <Img bg={b.bg} stripe={b.stripe} label={b.label} style={{ width:130, flexShrink:0 }} />
        <div style={{ flex:1, padding:'18px 20px' }}>
          <div style={{ display:'flex', justifyContent:'space-between', alignItems:'flex-start', marginBottom:8, gap:12 }}>
            <div>
              <div style={{ display:'flex', gap:8, marginBottom:6, flexWrap:'wrap' }}>
                <span className={'status ' + st.cls}>{st.label}</span>
                <span className="badge badge-neutral">{typeIcon[b.type]} {b.type === 'tour' ? 'Tour' : b.type === 'hotel' ? 'Khách sạn' : 'Combo'}</span>
              </div>
              <div style={{ fontWeight:700, fontSize:15, lineHeight:1.3 }}>{b.product}</div>
            </div>
            <div style={{ textAlign:'right', flexShrink:0 }}>
              <div style={{ fontSize:11, color:'var(--text3)', marginBottom:2 }}>Mã đặt chỗ</div>
              <div style={{ fontWeight:800, fontSize:14, color:'var(--primary)', letterSpacing:.5 }}>{b.id}</div>
            </div>
          </div>
          <div style={{ display:'flex', gap:16, fontSize:13, color:'var(--text2)', flexWrap:'wrap', marginBottom:12 }}>
            <span style={{ display:'flex', gap:4, alignItems:'center' }}><Ico.clock /> Ngày sử dụng: <strong style={{ color:'var(--text)' }}>{b.useDate}</strong></span>
            <span style={{ display:'flex', gap:4, alignItems:'center' }}><Ico.person /> {b.persons}</span>
            <span style={{ display:'flex', gap:4, alignItems:'center' }}><Ico.pin /> {b.agent}</span>
          </div>
          <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center' }}>
            <div>
              <div style={{ fontSize:12, color:'var(--text3)' }}>Đặt ngày {b.bookDate}</div>
              <div style={{ fontSize:18, fontWeight:900, color:'var(--accent)' }}>{formatPriceFull(b.total)}</div>
            </div>
            <div style={{ display:'flex', gap:8 }}>
              {b.status === 'confirmed' && <button className="btn btn-accent btn-sm">Xem thanh toán</button>}
              {b.status === 'completed' && <button className="btn btn-outline btn-sm">Viết đánh giá</button>}
              {b.status === 'paid' && <button className="btn btn-outline btn-sm">Xem chi tiết</button>}
              {b.status === 'cancelled' && <button onClick={() => navigate('tour-detail',{tour:TOURS[0]})} className="btn btn-ghost btn-sm">Đặt lại</button>}
              <button className="btn btn-ghost btn-sm">Chi tiết</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

function LoyaltyPanel() {
  const tiers = [
    { name:'Đồng', min:0, max:5, color:'#a16207', bg:'#fef9c3', perks:['Tích 1 điểm/10.000đ','Ưu đãi sinh nhật'] },
    { name:'Bạc', min:5, max:15, color:'#6b7280', bg:'#f3f4f6', perks:['Tích x1.2 điểm','Hỗ trợ ưu tiên','Sale sớm 24h'] },
    { name:'Vàng', min:15, max:50, color:'#92400e', bg:'#fef3c7', perks:['Tích x1.5 điểm','Hỗ trợ ưu tiên','Voucher sinh nhật 500k'] },
    { name:'Bạch kim', min:50, max:999, color:'#1e3a5f', bg:'#e0effe', perks:['Tích x2 điểm','Hỗ trợ VIP 24/7','Quà sinh nhật đặc biệt','Nâng hạng phòng miễn phí'] },
  ];
  const cur = tiers.find(t => USER.totalSpend/1e6 >= t.min && USER.totalSpend/1e6 < t.max) || tiers[2];
  const next = tiers[tiers.indexOf(cur)+1];
  const prog = next ? ((USER.totalSpend/1e6 - cur.min) / (next.min - cur.min) * 100) : 100;
  return (
    <div style={{ display:'flex', flexDirection:'column', gap:20 }}>
      {/* Points card */}
      <div style={{ background:'linear-gradient(135deg, var(--primary) 0%, oklch(0.28 0.16 245) 100%)', borderRadius:'var(--r-xl)', padding:'28px', color:'#fff' }}>
        <div style={{ display:'flex', justifyContent:'space-between', alignItems:'flex-start', marginBottom:20 }}>
          <div>
            <div style={{ fontSize:13, fontWeight:600, color:'rgba(255,255,255,.7)', marginBottom:8, textTransform:'uppercase', letterSpacing:1 }}>Điểm tích lũy</div>
            <div style={{ fontSize:52, fontWeight:900, lineHeight:1 }}>{USER.points.toLocaleString()}</div>
            <div style={{ fontSize:14, color:'rgba(255,255,255,.7)', marginTop:6 }}>≈ {formatPriceFull(USER.points*100)} · Hết hạn 01/03/2027</div>
          </div>
          <TierBadge tier={USER.tier} />
        </div>
        {next && <div>
          <div style={{ display:'flex', justifyContent:'space-between', fontSize:13, color:'rgba(255,255,255,.7)', marginBottom:6 }}>
            <span>Tiến tới hạng <strong style={{ color:'#fff' }}>{next.name}</strong></span>
            <span>{formatPrice((next.min - USER.totalSpend/1e6)*1e6)} nữa</span>
          </div>
          <div style={{ height:8, background:'rgba(255,255,255,.2)', borderRadius:4, overflow:'hidden' }}>
            <div style={{ height:'100%', background:'rgba(255,255,255,.85)', borderRadius:4, width:Math.min(100,prog)+'%', transition:'width .6s' }} />
          </div>
        </div>}
      </div>

      {/* Tier cards */}
      <div style={{ fontWeight:800, fontSize:18, marginBottom:4 }}>Hạng thành viên</div>
      <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:14 }}>
        {tiers.map(t => {
          const active = t.name === USER.tier;
          return (
            <div key={t.name} style={{ padding:'18px', borderRadius:'var(--r-lg)', background: active ? t.bg : '#fff', border: active ? '2px solid ' + t.color : '1.5px solid var(--border-lt)', transition:'all .2s' }}>
              <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center', marginBottom:10 }}>
                <span style={{ fontWeight:800, fontSize:16, color: t.color }}>{t.name}</span>
                {active && <span style={{ fontSize:11, fontWeight:700, background: t.color, color:'#fff', padding:'2px 8px', borderRadius:20 }}>Hạng hiện tại</span>}
              </div>
              <div style={{ fontSize:12, color:'var(--text3)', marginBottom:10 }}>Chi tiêu {t.min}–{t.max < 999 ? t.max : '50+'}tr/năm</div>
              <ul style={{ listStyle:'none', display:'flex', flexDirection:'column', gap:5 }}>
                {t.perks.map(p => <li key={p} style={{ fontSize:13, color:'var(--text2)', display:'flex', gap:6, alignItems:'center' }}><span style={{ color:t.color }}><Ico.check /></span>{p}</li>)}
              </ul>
            </div>
          );
        })}
      </div>

      {/* Redeem */}
      <div style={{ background:'#fff', borderRadius:'var(--r-lg)', padding:'22px', border:'1px solid var(--border-lt)' }}>
        <div style={{ fontWeight:700, fontSize:16, marginBottom:14 }}>Đổi điểm lấy ưu đãi</div>
        <div style={{ display:'grid', gridTemplateColumns:'repeat(3,1fr)', gap:12 }}>
          {[[100,'10.000đ giảm giá'],[500,'50.000đ giảm giá'],[1000,'100.000đ giảm giá']].map(([pts,reward]) => (
            <div key={pts} style={{ padding:'14px', borderRadius:'var(--r-md)', border:'1.5px solid var(--border)', textAlign:'center' }}>
              <div style={{ fontSize:18, fontWeight:900, color:'var(--primary)' }}>{pts.toLocaleString()} điểm</div>
              <div style={{ fontSize:13, color:'var(--text2)', margin:'4px 0 10px' }}>{reward}</div>
              <button disabled={USER.points < pts} className="btn btn-outline btn-sm" style={{ width:'100%' }}>Đổi ngay</button>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

function VouchersPanel() {
  const vouchers = [
    { code:'SUMMER30', disc:'30%', exp:'30/06/2026', min:3000000, cond:'Áp dụng cho tour & combo', used:false, type:'percent' },
    { code:'HALONG200K', disc:'200.000đ', exp:'15/06/2026', min:2000000, cond:'Chỉ dành cho tour Hạ Long', used:false, type:'fixed' },
    { code:'BIRTHDAY500', disc:'500.000đ', exp:'31/05/2026', min:5000000, cond:'Quà sinh nhật hạng Vàng', used:true, type:'fixed' },
  ];
  return (
    <div>
      <div style={{ background:'#fff', borderRadius:'var(--r-lg)', padding:'18px 22px', border:'1px solid var(--border-lt)', marginBottom:20, display:'flex', alignItems:'center', gap:12 }}>
        <h2 style={{ fontWeight:800, fontSize:18, flex:1 }}>Voucher của tôi</h2>
        <span className="badge badge-accent">{vouchers.filter(v=>!v.used).length} voucher khả dụng</span>
      </div>
      <div style={{ display:'flex', flexDirection:'column', gap:14 }}>
        {vouchers.map(v => (
          <div key={v.code} style={{ background:'#fff', borderRadius:'var(--r-lg)', border:'1.5px solid var(--border-lt)', overflow:'hidden', opacity: v.used ? 0.6 : 1, display:'flex' }}>
            <div style={{ width:8, background: v.used ? 'var(--border)' : 'var(--accent)', flexShrink:0 }} />
            <div style={{ flex:1, padding:'18px 20px', display:'flex', gap:16, alignItems:'center' }}>
              <div style={{ width:64, height:64, borderRadius:'var(--r-md)', background: v.used ? 'var(--border-lt)' : 'var(--accent-lt)', display:'flex', flexDirection:'column', alignItems:'center', justifyContent:'center', flexShrink:0 }}>
                <Ico.percent />
                <span style={{ fontSize:13, fontWeight:900, color: v.used ? 'var(--text3)' : 'var(--accent)', marginTop:2 }}>{v.disc}</span>
              </div>
              <div style={{ flex:1 }}>
                <div style={{ display:'flex', gap:8, marginBottom:4 }}>
                  <span style={{ fontWeight:800, fontSize:16, letterSpacing:1, color:'var(--text)', fontFamily:'monospace' }}>{v.code}</span>
                  {v.used && <span className="badge badge-neutral">Đã dùng</span>}
                </div>
                <div style={{ fontSize:13, color:'var(--text2)', marginBottom:2 }}>{v.cond}</div>
                <div style={{ fontSize:12, color:'var(--text3)' }}>Đơn tối thiểu {formatPrice(v.min)} · HSD: {v.exp}</div>
              </div>
              {!v.used && <button className="btn btn-accent btn-sm" style={{ flexShrink:0 }}>Dùng ngay</button>}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

function ReviewsPanel() {
  return (
    <div style={{ display:'flex', flexDirection:'column', gap:14 }}>
      <div style={{ background:'#fff', borderRadius:'var(--r-lg)', padding:'18px 22px', border:'1px solid var(--border-lt)', display:'flex', justifyContent:'space-between', alignItems:'center' }}>
        <h2 style={{ fontWeight:800, fontSize:18 }}>Đánh giá của tôi</h2>
        <span style={{ fontSize:13, color:'var(--text3)' }}>3 đánh giá</span>
      </div>
      {BOOKINGS.filter(b=>b.status==='completed').map(b => (
        <div key={b.id} style={{ background:'#fff', borderRadius:'var(--r-lg)', padding:'20px', border:'1px solid var(--border-lt)', display:'flex', gap:14 }}>
          <Img bg={b.bg} stripe={b.stripe} label="" style={{ width:80, height:60, borderRadius:'var(--r-sm)', flexShrink:0 }} />
          <div style={{ flex:1 }}>
            <div style={{ fontWeight:700, fontSize:14, marginBottom:6 }}>{b.product}</div>
            <div style={{ fontSize:12, color:'var(--text3)', marginBottom:10 }}>Sử dụng ngày {b.useDate} · {b.agent}</div>
            <div style={{ display:'flex', gap:2, marginBottom:8 }}>
              {[1,2,3,4,5].map(s => (
                <button key={s} style={{ background:'none', border:'none', cursor:'pointer', fontSize:22, color:'var(--gold)', padding:'0 2px' }}>★</button>
              ))}
            </div>
            <textarea className="input" rows={2} placeholder="Chia sẻ trải nghiệm của bạn..." style={{ fontSize:13 }} />
            <button className="btn btn-primary btn-sm" style={{ marginTop:8 }}>Gửi đánh giá</button>
          </div>
        </div>
      ))}
    </div>
  );
}

function AffiliatePanel() {
  return (
    <div style={{ display:'flex', flexDirection:'column', gap:20 }}>
      <div style={{ background:'linear-gradient(135deg, var(--accent) 0%, oklch(0.52 0.22 20) 100%)', borderRadius:'var(--r-xl)', padding:'28px', color:'#fff' }}>
        <div style={{ fontSize:13, fontWeight:700, color:'rgba(255,255,255,.7)', marginBottom:8, textTransform:'uppercase', letterSpacing:1 }}>Chương trình cộng tác viên</div>
        <div style={{ fontSize:30, fontWeight:900, marginBottom:8 }}>Kiếm hoa hồng từ mỗi booking</div>
        <p style={{ fontSize:14, color:'rgba(255,255,255,.8)', marginBottom:20 }}>Chia sẻ link tour & khách sạn, nhận hoa hồng 3–8% cho mỗi booking thành công. Thanh toán hàng tháng.</p>
        <button className="btn" style={{ background:'#fff', color:'var(--accent)', fontWeight:800 }}>Đăng ký ngay</button>
      </div>
      <div style={{ display:'grid', gridTemplateColumns:'repeat(3,1fr)', gap:14 }}>
        {[['0','Tổng click'],['0','Lượt chuyển đổi'],['0đ','Hoa hồng tích lũy']].map(([val,lbl]) => (
          <div key={lbl} style={{ background:'#fff', borderRadius:'var(--r-lg)', padding:'20px', border:'1px solid var(--border-lt)', textAlign:'center' }}>
            <div style={{ fontSize:28, fontWeight:900, color:'var(--primary)' }}>{val}</div>
            <div style={{ fontSize:13, color:'var(--text3)', marginTop:4 }}>{lbl}</div>
          </div>
        ))}
      </div>
    </div>
  );
}

function ProfilePanel() {
  return (
    <div style={{ background:'#fff', borderRadius:'var(--r-lg)', padding:'28px', border:'1px solid var(--border-lt)' }}>
      <h2 style={{ fontWeight:800, fontSize:18, marginBottom:24 }}>Thông tin cá nhân</h2>
      <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:18 }}>
        {[['Họ và tên', USER.name],['Số điện thoại', USER.phone],['Email', USER.email],['Ngày sinh','15/08/1990'],['Giới tính','Nam'],['Địa chỉ','Hà Nội, Việt Nam']].map(([lbl,val]) => (
          <div key={lbl}>
            <label style={{ fontSize:12, fontWeight:700, color:'var(--text3)', letterSpacing:.5, textTransform:'uppercase', display:'block', marginBottom:6 }}>{lbl}</label>
            <input className="input" defaultValue={val} style={{ fontSize:14 }} />
          </div>
        ))}
      </div>
      <button className="btn btn-primary" style={{ marginTop:20 }}>Lưu thay đổi</button>
    </div>
  );
}

function SecurityPanel() {
  return (
    <div style={{ display:'flex', flexDirection:'column', gap:14 }}>
      <div style={{ background:'#fff', borderRadius:'var(--r-lg)', padding:'22px', border:'1px solid var(--border-lt)' }}>
        <h2 style={{ fontWeight:800, fontSize:18, marginBottom:20 }}>Đổi mật khẩu</h2>
        {['Mật khẩu hiện tại','Mật khẩu mới','Xác nhận mật khẩu mới'].map(lbl => (
          <div key={lbl} style={{ marginBottom:14 }}>
            <label style={{ fontSize:13, fontWeight:600, display:'block', marginBottom:6 }}>{lbl}</label>
            <input type="password" className="input" placeholder="••••••••" style={{ fontSize:14 }} />
          </div>
        ))}
        <button className="btn btn-primary">Đổi mật khẩu</button>
      </div>
      <div style={{ background:'#fff', borderRadius:'var(--r-lg)', padding:'22px', border:'1px solid var(--border-lt)' }}>
        <h2 style={{ fontWeight:800, fontSize:18, marginBottom:4 }}>Xác thực 2 lớp</h2>
        <p style={{ fontSize:13, color:'var(--text2)', marginBottom:14 }}>Bảo vệ tài khoản bằng mã OTP khi đăng nhập.</p>
        <button className="btn btn-outline">Bật xác thực 2 lớp</button>
      </div>
    </div>
  );
}

Object.assign(window, { AccountPage });
