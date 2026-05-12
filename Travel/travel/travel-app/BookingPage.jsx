// BookingPage — 3-step booking form
function BookingPage({ navigate, params }) {
  const tour = params?.tour || TOURS[0];
  const [step, setStep] = useState(1);
  const [adults, setAdults] = useState(params?.adults || 2);
  const [children, setChildren] = useState(params?.children || 0);
  const [selDep, setSelDep] = useState(DEPARTURES[0]);
  const [contact, setContact] = useState({ name: USER.name, phone: USER.phone, email: USER.email, note: '' });
  const [passengers, setPassengers] = useState([
    { name: USER.name, dob: '', gender: 'male', cccd: '' },
    { name: '', dob: '', gender: 'male', cccd: '' },
  ]);
  const [voucher, setVoucher] = useState('');
  const [vOk, setVOk] = useState(false);
  const [usePoints, setUsePoints] = useState(false);
  const [agreed, setAgreed] = useState(false);
  const [done, setDone] = useState(false);

  const priceAdult = selDep ? selDep.priceAdult : tour.price;
  const priceChild = selDep ? selDep.priceChild : Math.round(tour.price * 0.7);
  const subtotal = adults * priceAdult + children * priceChild;
  const vDisc = vOk ? Math.round(subtotal * 0.1) : 0;
  const ptDisc = usePoints ? Math.min(USER.points * 100, Math.round(subtotal * 0.05)) : 0;
  const total = subtotal - vDisc - ptDisc;

  const stepsCfg = [
    { n:1, label:'Chọn ngày & hành khách' },
    { n:2, label:'Thông tin đặt chỗ' },
    { n:3, label:'Xác nhận & gửi' },
  ];

  function updatePassenger(i, k, v) {
    setPassengers(p => p.map((x, j) => j === i ? { ...x, [k]: v } : x));
  }

  if (done) return <BookingSuccess navigate={navigate} bookId="BK260512" total={total} tour={tour} />;

  return (
    <div className="page-anim" style={{ paddingTop:68, background:'var(--bg)', minHeight:'100vh' }}>
      <div className="container" style={{ paddingTop:32, paddingBottom:48 }}>
        <button onClick={() => step > 1 ? setStep(s => s-1) : navigate('tour-detail',{tour})}
          style={{ background:'none', border:'none', cursor:'pointer', color:'var(--primary)', fontWeight:600, fontSize:14, display:'flex', alignItems:'center', gap:6, marginBottom:24, padding:0 }}>
          <Ico.chevL /> Quay lại
        </button>

        {/* Progress */}
        <div style={{ display:'flex', alignItems:'center', maxWidth:600, marginBottom:36 }}>
          {stepsCfg.map((s, i) => (
            <React.Fragment key={s.n}>
              <div style={{ display:'flex', flexDirection:'column', alignItems:'center', gap:6 }}>
                <div className="step-dot"
                  style={{ background: step > s.n ? 'var(--success)' : step === s.n ? 'var(--primary)' : 'var(--border)', color: step >= s.n ? '#fff' : 'var(--text3)' }}>
                  {step > s.n ? <Ico.check /> : s.n}
                </div>
                <span style={{ fontSize:12, fontWeight: step === s.n ? 700 : 400, color: step === s.n ? 'var(--primary)' : 'var(--text3)', whiteSpace:'nowrap' }}>{s.label}</span>
              </div>
              {i < stepsCfg.length - 1 && <div className={'step-line' + (step > s.n ? ' done' : '')} style={{ marginBottom:20 }} />}
            </React.Fragment>
          ))}
        </div>

        <div style={{ display:'grid', gridTemplateColumns:'1fr 340px', gap:24, alignItems:'flex-start' }}>
          {/* ── Main form ── */}
          <div>

            {/* ── STEP 1 ── */}
            {step === 1 && (
              <div style={{ display:'flex', flexDirection:'column', gap:20 }}>
                <FormCard title="Chọn đợt khởi hành">
                  <div style={{ display:'flex', flexDirection:'column', gap:8 }}>
                    {DEPARTURES.filter(d => d.booked < d.slots).map(d => (
                      <button key={d.id} onClick={() => setSelDep(d)}
                        style={{ padding:'14px 18px', borderRadius:'var(--r-md)', border: selDep?.id===d.id ? '2px solid var(--primary)' : '1.5px solid var(--border)', background: selDep?.id===d.id ? 'var(--primary-xlt)' : '#fff', cursor:'pointer', textAlign:'left', transition:'all .15s', display:'flex', justifyContent:'space-between', alignItems:'center' }}>
                        <div>
                          <span style={{ fontWeight:700, fontSize:15, color: selDep?.id===d.id ? 'var(--primary)' : 'var(--text)' }}>{d.date}</span>
                          <span style={{ fontSize:13, color:'var(--text3)', marginLeft:8 }}>{d.dayOfWeek}</span>
                          <div style={{ fontSize:13, color:'var(--text2)', marginTop:2 }}>Còn {d.slots-d.booked}/{d.slots} chỗ trống</div>
                        </div>
                        <div style={{ textAlign:'right' }}>
                          <div style={{ fontWeight:800, fontSize:16, color:'var(--accent)' }}>{formatPriceFull(d.priceAdult)}</div>
                          <div style={{ fontSize:12, color:'var(--text3)' }}>/người lớn</div>
                        </div>
                      </button>
                    ))}
                  </div>
                </FormCard>
                <FormCard title="Số hành khách">
                  <div style={{ display:'flex', flexDirection:'column', gap:16 }}>
                    <NumStepper label="Người lớn" sublabel="Từ 12 tuổi trở lên" value={adults} onChange={v => { setAdults(v); setPassengers(Array.from({length:v+children}, (_,i)=>passengers[i]||{name:'',dob:'',gender:'male',cccd:''})); }} min={1} max={20} />
                    <hr />
                    <NumStepper label="Trẻ em" sublabel="Từ 5–11 tuổi · 70% giá vé" value={children} onChange={v => { setChildren(v); setPassengers(Array.from({length:adults+v}, (_,i)=>passengers[i]||{name:'',dob:'',gender:'male',cccd:''})); }} min={0} max={10} />
                    <hr />
                    <div style={{ display:'flex', alignItems:'center', justifyContent:'space-between' }}>
                      <div>
                        <div style={{ fontSize:14, fontWeight:600 }}>Em bé</div>
                        <div style={{ fontSize:12, color:'var(--text3)' }}>Dưới 5 tuổi · Miễn phí · Không có ghế riêng</div>
                      </div>
                      <span className="badge badge-success">Miễn phí</span>
                    </div>
                  </div>
                </FormCard>
                <button onClick={() => setStep(2)} className="btn btn-primary btn-lg" style={{ alignSelf:'flex-start', gap:8 }}>
                  Tiếp tục <Ico.chevR />
                </button>
              </div>
            )}

            {/* ── STEP 2 ── */}
            {step === 2 && (
              <div style={{ display:'flex', flexDirection:'column', gap:20 }}>
                <FormCard title="Thông tin liên hệ">
                  <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:14 }}>
                    <FormField label="Họ và tên *" value={contact.name} onChange={v => setContact({...contact,name:v})} placeholder="Nguyễn Văn A" />
                    <FormField label="Số điện thoại *" value={contact.phone} onChange={v => setContact({...contact,phone:v})} placeholder="0901 234 567" />
                    <FormField label="Email *" value={contact.email} onChange={v => setContact({...contact,email:v})} placeholder="email@example.com" cls="colspan-2" />
                  </div>
                  <div style={{ marginTop:14 }}>
                    <label style={{ fontSize:13, fontWeight:600, display:'block', marginBottom:6 }}>Yêu cầu đặc biệt</label>
                    <textarea className="input" rows={3} value={contact.note} onChange={e => setContact({...contact,note:e.target.value})}
                      placeholder="Chế độ ăn đặc biệt, hỗ trợ người khuyết tật, phòng riêng biệt..." style={{ resize:'vertical', fontSize:13 }} />
                  </div>
                </FormCard>
                {Array.from({length: adults + children}, (_,i) => (
                  <FormCard key={i} title={`Hành khách ${i+1}${i < adults ? ' (Người lớn)' : ' (Trẻ em)'}`}>
                    <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:14 }}>
                      <FormField label="Họ và tên *" value={passengers[i]?.name||''} onChange={v => updatePassenger(i,'name',v)} placeholder="Nguyễn Văn A" />
                      <div>
                        <label style={{ fontSize:13, fontWeight:600, display:'block', marginBottom:6 }}>Giới tính *</label>
                        <div style={{ display:'flex', gap:10 }}>
                          {[['male','Nam'],['female','Nữ']].map(([val,lbl]) => (
                            <label key={val} style={{ display:'flex', gap:6, alignItems:'center', cursor:'pointer', fontSize:14 }}>
                              <input type="radio" name={'gender_'+i} value={val} checked={passengers[i]?.gender===val} onChange={() => updatePassenger(i,'gender',val)} style={{ accentColor:'var(--primary)' }} /> {lbl}
                            </label>
                          ))}
                        </div>
                      </div>
                      <FormField label="Ngày sinh *" value={passengers[i]?.dob||''} onChange={v => updatePassenger(i,'dob',v)} type="date" />
                      <FormField label="CCCD / Hộ chiếu" value={passengers[i]?.cccd||''} onChange={v => updatePassenger(i,'cccd',v)} placeholder="012345678901" />
                    </div>
                  </FormCard>
                ))}
                <button onClick={() => setStep(3)} className="btn btn-primary btn-lg" style={{ alignSelf:'flex-start', gap:8 }}>Tiếp tục <Ico.chevR /></button>
              </div>
            )}

            {/* ── STEP 3 ── */}
            {step === 3 && (
              <div style={{ display:'flex', flexDirection:'column', gap:20 }}>
                <FormCard title="Tóm tắt đặt chỗ">
                  <div style={{ display:'flex', gap:14, alignItems:'flex-start' }}>
                    <Img bg={tour.bg} stripe={tour.stripe} label="" style={{ width:100, height:72, borderRadius:'var(--r-sm)', flexShrink:0 }} />
                    <div style={{ flex:1 }}>
                      <div style={{ fontWeight:700, fontSize:15, marginBottom:4 }}>{tour.name}</div>
                      <div style={{ fontSize:13, color:'var(--text2)', display:'flex', gap:10, flexWrap:'wrap' }}>
                        <span style={{ display:'flex', gap:3, alignItems:'center' }}><Ico.clock /> {selDep?.date || ''}</span>
                        <span style={{ display:'flex', gap:3, alignItems:'center' }}><Ico.person /> {adults} NL{children>0?', '+children+' TE':''}</span>
                        <span style={{ display:'flex', gap:3, alignItems:'center' }}><Ico.pin /> {tour.departure} → {tour.dest.split(',')[0]}</span>
                      </div>
                    </div>
                  </div>
                </FormCard>
                <FormCard title="Mã giảm giá & điểm thưởng">
                  <div style={{ display:'flex', gap:8, marginBottom:14 }}>
                    <input className="input" placeholder="Nhập mã voucher..." value={voucher} onChange={e => setVoucher(e.target.value.toUpperCase())} style={{ flex:1, fontSize:13 }} />
                    <button onClick={() => voucher && setVOk(true)} className="btn btn-outline btn-sm">Áp dụng</button>
                  </div>
                  {vOk && <div style={{ fontSize:12, color:'var(--success)', fontWeight:600, marginBottom:12, display:'flex', gap:4, alignItems:'center' }}><Ico.check /> Giảm 10% · {voucher}</div>}
                  <label style={{ display:'flex', gap:10, alignItems:'center', cursor:'pointer', padding:'12px 14px', borderRadius:'var(--r-md)', border:'1.5px solid var(--border)', background: usePoints ? 'var(--gold-lt)' : '#fff' }}>
                    <input type="checkbox" checked={usePoints} onChange={e => setUsePoints(e.target.checked)} style={{ accentColor:'var(--primary)', width:16, height:16 }} />
                    <div style={{ flex:1 }}>
                      <div style={{ fontWeight:600, fontSize:14 }}>Dùng điểm thưởng</div>
                      <div style={{ fontSize:12, color:'var(--text2)' }}>Bạn có <strong>{USER.points.toLocaleString()}</strong> điểm (trị giá {formatPriceFull(USER.points*100)}) — Giảm tối đa 5%</div>
                    </div>
                    <TierBadge tier={USER.tier} />
                  </label>
                </FormCard>
                <FormCard title="Điều khoản & Xác nhận">
                  <div style={{ fontSize:13, color:'var(--text2)', lineHeight:1.7, marginBottom:14 }}>
                    Bằng cách đặt chỗ, bạn xác nhận đã đọc và đồng ý với <a href="#" style={{ color:'var(--primary)', fontWeight:600 }}>Điều khoản sử dụng</a>, <a href="#" style={{ color:'var(--primary)', fontWeight:600 }}>Chính sách bảo mật</a> và <a href="#" style={{ color:'var(--primary)', fontWeight:600 }}>Chính sách hủy tour</a> của TravelViet.
                  </div>
                  <label style={{ display:'flex', gap:10, alignItems:'flex-start', cursor:'pointer' }}>
                    <input type="checkbox" checked={agreed} onChange={e => setAgreed(e.target.checked)} style={{ accentColor:'var(--primary)', width:16, height:16, marginTop:2, flexShrink:0 }} />
                    <span style={{ fontSize:14 }}>Tôi đồng ý với các điều khoản và xác nhận thông tin hành khách là chính xác.</span>
                  </label>
                </FormCard>
                <div style={{ padding:'18px 20px', background:'var(--primary-xlt)', borderRadius:'var(--r-md)', border:'1px solid var(--primary-lt)', fontSize:13, color:'var(--text2)', lineHeight:1.6 }}>
                  <strong style={{ color:'var(--primary)' }}>Hướng dẫn thanh toán:</strong> Sau khi đặt chỗ thành công, đại lý sẽ xác nhận trong vòng 24h và gửi hướng dẫn thanh toán qua email. Vui lòng không chuyển tiền cho đến khi nhận được xác nhận.
                </div>
                <button onClick={() => agreed && setDone(true)} disabled={!agreed} className="btn btn-accent btn-lg" style={{ alignSelf:'flex-start', gap:8 }}>
                  Gửi yêu cầu đặt chỗ — {formatPrice(total)}
                </button>
              </div>
            )}
          </div>

          {/* ── Price sidebar ── */}
          <div style={{ position:'sticky', top:84 }}>
            <div className="card-flat" style={{ borderRadius:'var(--r-xl)', padding:'22px', border:'2px solid var(--border)' }}>
              <div style={{ fontWeight:800, fontSize:15, marginBottom:14 }}>Tóm tắt chi phí</div>
              <div style={{ display:'flex', gap:10, marginBottom:14 }}>
                <Img bg={tour.bg} stripe={tour.stripe} label="" style={{ width:72, height:54, borderRadius:'var(--r-sm)', flexShrink:0 }} />
                <div style={{ fontSize:13, fontWeight:600, lineHeight:1.4 }}>{tour.name}</div>
              </div>
              {selDep && <div style={{ fontSize:13, color:'var(--text2)', marginBottom:14, padding:'10px 12px', background:'var(--bg)', borderRadius:'var(--r-sm)' }}>
                Khởi hành: <strong>{selDep.date}</strong> · {selDep.dayOfWeek}
              </div>}
              <hr style={{ marginBottom:14 }} />
              <div style={{ display:'flex', flexDirection:'column', gap:8, marginBottom:14 }}>
                <PriceLine label={`NL × ${adults}`} val={adults * priceAdult} />
                {children > 0 && <PriceLine label={`TE × ${children}`} val={children * priceChild} />}
                {vOk && <PriceLine label="Voucher −10%" val={-vDisc} accent="var(--success)" />}
                {usePoints && ptDisc > 0 && <PriceLine label="Điểm thưởng" val={-ptDisc} accent="var(--success)" />}
              </div>
              <hr style={{ marginBottom:14 }} />
              <div style={{ display:'flex', justifyContent:'space-between', marginBottom:18 }}>
                <span style={{ fontWeight:800, fontSize:16 }}>Tổng cộng</span>
                <span style={{ fontWeight:900, fontSize:20, color:'var(--accent)' }}>{formatPrice(total)}</span>
              </div>
              <div style={{ padding:'10px 14px', background:'var(--success-lt)', borderRadius:'var(--r-md)', fontSize:12, color:'var(--success)', fontWeight:600, display:'flex', gap:6, alignItems:'center' }}>
                <Ico.check /> Thanh toán sau khi đại lý xác nhận
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

function FormCard({ title, children }) {
  return (
    <div className="card-flat" style={{ borderRadius:'var(--r-lg)', padding:'22px 24px' }}>
      <div style={{ fontWeight:800, fontSize:16, marginBottom:18, paddingBottom:14, borderBottom:'1px solid var(--border-lt)' }}>{title}</div>
      {children}
    </div>
  );
}

function FormField({ label, value, onChange, placeholder, type = 'text' }) {
  return (
    <div>
      <label style={{ fontSize:13, fontWeight:600, display:'block', marginBottom:6 }}>{label}</label>
      <input type={type} className="input" value={value} onChange={e => onChange(e.target.value)} placeholder={placeholder} style={{ fontSize:13 }} />
    </div>
  );
}

function PriceLine({ label, val, accent }) {
  return (
    <div style={{ display:'flex', justifyContent:'space-between', fontSize:13, color: accent || 'var(--text2)' }}>
      <span>{label}</span>
      <span style={{ fontWeight:600 }}>{val < 0 ? '−' : ''}{formatPriceFull(Math.abs(val))}</span>
    </div>
  );
}

function BookingSuccess({ navigate, bookId, total, tour }) {
  return (
    <div className="page-anim" style={{ paddingTop:68, minHeight:'100vh', display:'flex', alignItems:'center', justifyContent:'center', background:'var(--bg)' }}>
      <div style={{ maxWidth:540, textAlign:'center', padding:'0 24px' }}>
        <div style={{ width:80, height:80, borderRadius:'50%', background:'var(--success-lt)', display:'flex', alignItems:'center', justifyContent:'center', margin:'0 auto 24px', fontSize:36 }}>✅</div>
        <h2 style={{ fontSize:28, fontWeight:900, marginBottom:10 }}>Đặt chỗ thành công!</h2>
        <p style={{ fontSize:15, color:'var(--text2)', marginBottom:20 }}>Yêu cầu của bạn đã được gửi đến đại lý. Vui lòng chờ xác nhận trong vòng 24 giờ.</p>
        <div style={{ background:'#fff', borderRadius:'var(--r-xl)', padding:'24px', marginBottom:24, border:'2px solid var(--border)' }}>
          <div style={{ fontSize:13, color:'var(--text3)', marginBottom:4 }}>Mã đặt chỗ của bạn</div>
          <div style={{ fontSize:28, fontWeight:900, color:'var(--primary)', letterSpacing:2, marginBottom:12 }}>{bookId}</div>
          <div style={{ fontSize:14, color:'var(--text2)' }}>{tour.name}</div>
          <div style={{ fontSize:16, fontWeight:800, color:'var(--accent)', marginTop:4 }}>{formatPrice(total)}</div>
        </div>
        <div style={{ padding:'14px 18px', background:'var(--warn-lt)', borderRadius:'var(--r-md)', fontSize:13, color:'var(--text2)', marginBottom:24, textAlign:'left' }}>
          Chúng tôi đã gửi xác nhận đến <strong>{USER.email}</strong>. Đại lý sẽ liên hệ hướng dẫn thanh toán.
        </div>
        <div style={{ display:'flex', gap:12, justifyContent:'center' }}>
          <button onClick={() => navigate('account')} className="btn btn-primary btn-lg">Xem đặt chỗ của tôi</button>
          <button onClick={() => navigate('home')} className="btn btn-outline btn-lg">Về trang chủ</button>
        </div>
      </div>
    </div>
  );
}

Object.assign(window, { BookingPage });
