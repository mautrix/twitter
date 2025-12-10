package com.x.dmv2.thriftjava;

import android.gov.nist.core.C0003b;
import android.gov.nist.core.Separators;
import com.bendb.thrifty.InterfaceC11261a;
import com.bendb.thrifty.kotlin.InterfaceC11262a;
import com.bendb.thrifty.protocol.C11265c;
import com.bendb.thrifty.protocol.InterfaceC11268f;
import com.bendb.thrifty.util.C11272a;
import java.io.IOException;
import kotlin.Metadata;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.Intrinsics;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

@Metadata(m64929d1 = {"\u0000:\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010\t\n\u0000\n\u0002\u0010\u000e\n\u0002\b\u0004\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\n\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\u000b\n\u0002\b\u0007\b\u0086\b\u0018\u0000  2\u00020\u0001:\u0002! B%\u0012\b\u0010\u0003\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0005\u001a\u0004\u0018\u00010\u0004\u0012\b\u0010\u0006\u001a\u0004\u0018\u00010\u0004¢\u0006\u0004\b\u0007\u0010\bJ\u0017\u0010\f\u001a\u00020\u000b2\u0006\u0010\n\u001a\u00020\tH\u0016¢\u0006\u0004\b\f\u0010\rJ\u0012\u0010\u000e\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u000e\u0010\u000fJ\u0012\u0010\u0010\u001a\u0004\u0018\u00010\u0004HÆ\u0003¢\u0006\u0004\b\u0010\u0010\u0011J\u0012\u0010\u0012\u001a\u0004\u0018\u00010\u0004HÆ\u0003¢\u0006\u0004\b\u0012\u0010\u0011J4\u0010\u0013\u001a\u00020\u00002\n\b\u0002\u0010\u0003\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0005\u001a\u0004\u0018\u00010\u00042\n\b\u0002\u0010\u0006\u001a\u0004\u0018\u00010\u0004HÆ\u0001¢\u0006\u0004\b\u0013\u0010\u0014J\u0010\u0010\u0015\u001a\u00020\u0004HÖ\u0001¢\u0006\u0004\b\u0015\u0010\u0011J\u0010\u0010\u0017\u001a\u00020\u0016HÖ\u0001¢\u0006\u0004\b\u0017\u0010\u0018J\u001a\u0010\u001c\u001a\u00020\u001b2\b\u0010\u001a\u001a\u0004\u0018\u00010\u0019HÖ\u0003¢\u0006\u0004\b\u001c\u0010\u001dR\u0016\u0010\u0003\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0003\u0010\u001eR\u0016\u0010\u0005\u001a\u0004\u0018\u00010\u00048\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0005\u0010\u001fR\u0016\u0010\u0006\u001a\u0004\u0018\u00010\u00048\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0006\u0010\u001f¨\u0006\""}, m64930d2 = {"Lcom/x/dmv2/thriftjava/GroupInviteEnable;", "Lcom/bendb/thrifty/a;", "", "expires_at_msec", "", "invite_url", "affiliate_id", "<init>", "(Ljava/lang/Long;Ljava/lang/String;Ljava/lang/String;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/lang/Long;", "component2", "()Ljava/lang/String;", "component3", "copy", "(Ljava/lang/Long;Ljava/lang/String;Ljava/lang/String;)Lcom/x/dmv2/thriftjava/GroupInviteEnable;", "toString", "", "hashCode", "()I", "", "other", "", "equals", "(Ljava/lang/Object;)Z", "Ljava/lang/Long;", "Ljava/lang/String;", "Companion", "GroupInviteEnableAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class GroupInviteEnable implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final String affiliate_id;

    @JvmField
    @InterfaceC88465b
    public final Long expires_at_msec;

    @JvmField
    @InterfaceC88465b
    public final String invite_url;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new GroupInviteEnableAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/GroupInviteEnable$GroupInviteEnableAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/GroupInviteEnable;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/GroupInviteEnable;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/GroupInviteEnable;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class GroupInviteEnableAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public GroupInviteEnable m85647read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Long lValueOf = null;
            String string = null;
            String string2 = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    return new GroupInviteEnable(lValueOf, string, string2);
                }
                short s = c11265cMo14127V2.f38393b;
                if (s != 1) {
                    if (s != 2) {
                        if (s != 3) {
                            C11272a.m14141a(protocol, b);
                        } else if (b == 11) {
                            string2 = protocol.readString();
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    } else if (b == 11) {
                        string = protocol.readString();
                    } else {
                        C11272a.m14141a(protocol, b);
                    }
                } else if (b == 10) {
                    lValueOf = Long.valueOf(protocol.mo14124H0());
                } else {
                    C11272a.m14141a(protocol, b);
                }
            }
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a GroupInviteEnable struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("GroupInviteEnable");
            if (struct.expires_at_msec != null) {
                protocol.mo14136v3("expires_at_msec", 1, (byte) 10);
                protocol.mo14121B3(struct.expires_at_msec.longValue());
            }
            if (struct.invite_url != null) {
                protocol.mo14136v3("invite_url", 2, (byte) 11);
                protocol.mo14137w0(struct.invite_url);
            }
            if (struct.affiliate_id != null) {
                protocol.mo14136v3("affiliate_id", 3, (byte) 11);
                protocol.mo14137w0(struct.affiliate_id);
            }
            protocol.mo14134i0();
        }
    }

    public GroupInviteEnable(@InterfaceC88465b Long l, @InterfaceC88465b String str, @InterfaceC88465b String str2) {
        this.expires_at_msec = l;
        this.invite_url = str;
        this.affiliate_id = str2;
    }

    public static /* synthetic */ GroupInviteEnable copy$default(GroupInviteEnable groupInviteEnable, Long l, String str, String str2, int i, Object obj) {
        if ((i & 1) != 0) {
            l = groupInviteEnable.expires_at_msec;
        }
        if ((i & 2) != 0) {
            str = groupInviteEnable.invite_url;
        }
        if ((i & 4) != 0) {
            str2 = groupInviteEnable.affiliate_id;
        }
        return groupInviteEnable.copy(l, str, str2);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final Long getExpires_at_msec() {
        return this.expires_at_msec;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final String getInvite_url() {
        return this.invite_url;
    }

    @InterfaceC88465b
    /* renamed from: component3, reason: from getter */
    public final String getAffiliate_id() {
        return this.affiliate_id;
    }

    @InterfaceC88464a
    public final GroupInviteEnable copy(@InterfaceC88465b Long expires_at_msec, @InterfaceC88465b String invite_url, @InterfaceC88465b String affiliate_id) {
        return new GroupInviteEnable(expires_at_msec, invite_url, affiliate_id);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof GroupInviteEnable)) {
            return false;
        }
        GroupInviteEnable groupInviteEnable = (GroupInviteEnable) other;
        return Intrinsics.m65267c(this.expires_at_msec, groupInviteEnable.expires_at_msec) && Intrinsics.m65267c(this.invite_url, groupInviteEnable.invite_url) && Intrinsics.m65267c(this.affiliate_id, groupInviteEnable.affiliate_id);
    }

    public int hashCode() {
        Long l = this.expires_at_msec;
        int iHashCode = (l == null ? 0 : l.hashCode()) * 31;
        String str = this.invite_url;
        int iHashCode2 = (iHashCode + (str == null ? 0 : str.hashCode())) * 31;
        String str2 = this.affiliate_id;
        return iHashCode2 + (str2 != null ? str2.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        Long l = this.expires_at_msec;
        String str = this.invite_url;
        String str2 = this.affiliate_id;
        StringBuilder sb = new StringBuilder("GroupInviteEnable(expires_at_msec=");
        sb.append(l);
        sb.append(", invite_url=");
        sb.append(str);
        sb.append(", affiliate_id=");
        return C0003b.m4b(sb, str2, Separators.RPAREN);
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}