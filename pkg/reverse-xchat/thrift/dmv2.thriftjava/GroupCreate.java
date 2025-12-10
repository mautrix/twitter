package com.x.dmv2.thriftjava;

import android.gov.nist.core.C0003b;
import android.gov.nist.core.Separators;
import android.gov.nist.javax.sip.clientauthutils.C0026b;
import com.bendb.thrifty.InterfaceC11261a;
import com.bendb.thrifty.kotlin.InterfaceC11262a;
import com.bendb.thrifty.protocol.C11265c;
import com.bendb.thrifty.protocol.InterfaceC11268f;
import com.bendb.thrifty.util.C11272a;
import java.io.IOException;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;
import kotlin.Metadata;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.Intrinsics;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

@Metadata(m64929d1 = {"\u00008\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010 \n\u0002\u0010\u000e\n\u0002\b\u0007\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\f\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\u000b\n\u0002\b\u0007\b\u0086\b\u0018\u0000 $2\u00020\u0001:\u0002%$BE\u0012\u000e\u0010\u0004\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u0002\u0012\u000e\u0010\u0005\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u0002\u0012\b\u0010\u0006\u001a\u0004\u0018\u00010\u0003\u0012\b\u0010\u0007\u001a\u0004\u0018\u00010\u0003\u0012\b\u0010\b\u001a\u0004\u0018\u00010\u0003¢\u0006\u0004\b\t\u0010\nJ\u0017\u0010\u000e\u001a\u00020\r2\u0006\u0010\f\u001a\u00020\u000bH\u0016¢\u0006\u0004\b\u000e\u0010\u000fJ\u0018\u0010\u0010\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u0010\u0010\u0011J\u0018\u0010\u0012\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u0012\u0010\u0011J\u0012\u0010\u0013\u001a\u0004\u0018\u00010\u0003HÆ\u0003¢\u0006\u0004\b\u0013\u0010\u0014J\u0012\u0010\u0015\u001a\u0004\u0018\u00010\u0003HÆ\u0003¢\u0006\u0004\b\u0015\u0010\u0014J\u0012\u0010\u0016\u001a\u0004\u0018\u00010\u0003HÆ\u0003¢\u0006\u0004\b\u0016\u0010\u0014JX\u0010\u0017\u001a\u00020\u00002\u0010\b\u0002\u0010\u0004\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u00022\u0010\b\u0002\u0010\u0005\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u00022\n\b\u0002\u0010\u0006\u001a\u0004\u0018\u00010\u00032\n\b\u0002\u0010\u0007\u001a\u0004\u0018\u00010\u00032\n\b\u0002\u0010\b\u001a\u0004\u0018\u00010\u0003HÆ\u0001¢\u0006\u0004\b\u0017\u0010\u0018J\u0010\u0010\u0019\u001a\u00020\u0003HÖ\u0001¢\u0006\u0004\b\u0019\u0010\u0014J\u0010\u0010\u001b\u001a\u00020\u001aHÖ\u0001¢\u0006\u0004\b\u001b\u0010\u001cJ\u001a\u0010 \u001a\u00020\u001f2\b\u0010\u001e\u001a\u0004\u0018\u00010\u001dHÖ\u0003¢\u0006\u0004\b \u0010!R\u001c\u0010\u0004\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0004\u0010\"R\u001c\u0010\u0005\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0005\u0010\"R\u0016\u0010\u0006\u001a\u0004\u0018\u00010\u00038\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0006\u0010#R\u0016\u0010\u0007\u001a\u0004\u0018\u00010\u00038\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0007\u0010#R\u0016\u0010\b\u001a\u0004\u0018\u00010\u00038\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\b\u0010#¨\u0006&"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/GroupCreate;", "Lcom/bendb/thrifty/a;", "", "", "member_ids", "admin_ids", "title", "avatar_url", "conversation_key_version", "<init>", "(Ljava/util/List;Ljava/util/List;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/util/List;", "component2", "component3", "()Ljava/lang/String;", "component4", "component5", "copy", "(Ljava/util/List;Ljava/util/List;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)Lcom/x/dmv2/thriftjava/GroupCreate;", "toString", "", "hashCode", "()I", "", "other", "", "equals", "(Ljava/lang/Object;)Z", "Ljava/util/List;", "Ljava/lang/String;", "Companion", "GroupCreateAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class GroupCreate implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final List admin_ids;

    @JvmField
    @InterfaceC88465b
    public final String avatar_url;

    @JvmField
    @InterfaceC88465b
    public final String conversation_key_version;

    @JvmField
    @InterfaceC88465b
    public final List member_ids;

    @JvmField
    @InterfaceC88465b
    public final String title;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new GroupCreateAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/GroupCreate$GroupCreateAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/GroupCreate;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/GroupCreate;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/GroupCreate;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class GroupCreateAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public GroupCreate m83718read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            ArrayList arrayList = null;
            ArrayList arrayList2 = null;
            String string = null;
            String string2 = null;
            String string3 = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    return new GroupCreate(arrayList, arrayList2, string, string2, string3);
                }
                int i = 0;
                short s = c11265cMo14127V2.f38393b;
                if (s != 1) {
                    if (s != 2) {
                        if (s != 3) {
                            if (s != 4) {
                                if (s != 5) {
                                    C11272a.m14141a(protocol, b);
                                } else if (b == 11) {
                                    string3 = protocol.readString();
                                } else {
                                    C11272a.m14141a(protocol, b);
                                }
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
                    } else if (b == 15) {
                        int i2 = protocol.mo14130a2().f38395b;
                        ArrayList arrayList3 = new ArrayList(i2);
                        while (i < i2) {
                            arrayList3.add(protocol.readString());
                            i++;
                        }
                        arrayList2 = arrayList3;
                    } else {
                        C11272a.m14141a(protocol, b);
                    }
                } else if (b == 15) {
                    int i3 = protocol.mo14130a2().f38395b;
                    ArrayList arrayList4 = new ArrayList(i3);
                    while (i < i3) {
                        arrayList4.add(protocol.readString());
                        i++;
                    }
                    arrayList = arrayList4;
                } else {
                    C11272a.m14141a(protocol, b);
                }
            }
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a GroupCreate struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("GroupCreate");
            if (struct.member_ids != null) {
                protocol.mo14136v3("member_ids", 1, (byte) 15);
                protocol.mo14128X0((byte) 11, struct.member_ids.size());
                Iterator it = struct.member_ids.iterator();
                while (it.hasNext()) {
                    protocol.mo14137w0((String) it.next());
                }
            }
            if (struct.admin_ids != null) {
                protocol.mo14136v3("admin_ids", 2, (byte) 15);
                protocol.mo14128X0((byte) 11, struct.admin_ids.size());
                Iterator it2 = struct.admin_ids.iterator();
                while (it2.hasNext()) {
                    protocol.mo14137w0((String) it2.next());
                }
            }
            if (struct.title != null) {
                protocol.mo14136v3("title", 3, (byte) 11);
                protocol.mo14137w0(struct.title);
            }
            if (struct.avatar_url != null) {
                protocol.mo14136v3("avatar_url", 4, (byte) 11);
                protocol.mo14137w0(struct.avatar_url);
            }
            if (struct.conversation_key_version != null) {
                protocol.mo14136v3("conversation_key_version", 5, (byte) 11);
                protocol.mo14137w0(struct.conversation_key_version);
            }
            protocol.mo14134i0();
        }
    }

    public GroupCreate(@InterfaceC88465b List list, @InterfaceC88465b List list2, @InterfaceC88465b String str, @InterfaceC88465b String str2, @InterfaceC88465b String str3) {
        this.member_ids = list;
        this.admin_ids = list2;
        this.title = str;
        this.avatar_url = str2;
        this.conversation_key_version = str3;
    }

    public static /* synthetic */ GroupCreate copy$default(GroupCreate groupCreate, List list, List list2, String str, String str2, String str3, int i, Object obj) {
        if ((i & 1) != 0) {
            list = groupCreate.member_ids;
        }
        if ((i & 2) != 0) {
            list2 = groupCreate.admin_ids;
        }
        List list3 = list2;
        if ((i & 4) != 0) {
            str = groupCreate.title;
        }
        String str4 = str;
        if ((i & 8) != 0) {
            str2 = groupCreate.avatar_url;
        }
        String str5 = str2;
        if ((i & 16) != 0) {
            str3 = groupCreate.conversation_key_version;
        }
        return groupCreate.copy(list, list3, str4, str5, str3);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final List getMember_ids() {
        return this.member_ids;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final List getAdmin_ids() {
        return this.admin_ids;
    }

    @InterfaceC88465b
    /* renamed from: component3, reason: from getter */
    public final String getTitle() {
        return this.title;
    }

    @InterfaceC88465b
    /* renamed from: component4, reason: from getter */
    public final String getAvatar_url() {
        return this.avatar_url;
    }

    @InterfaceC88465b
    /* renamed from: component5, reason: from getter */
    public final String getConversation_key_version() {
        return this.conversation_key_version;
    }

    @InterfaceC88464a
    public final GroupCreate copy(@InterfaceC88465b List member_ids, @InterfaceC88465b List admin_ids, @InterfaceC88465b String title, @InterfaceC88465b String avatar_url, @InterfaceC88465b String conversation_key_version) {
        return new GroupCreate(member_ids, admin_ids, title, avatar_url, conversation_key_version);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof GroupCreate)) {
            return false;
        }
        GroupCreate groupCreate = (GroupCreate) other;
        return Intrinsics.m65267c(this.member_ids, groupCreate.member_ids) && Intrinsics.m65267c(this.admin_ids, groupCreate.admin_ids) && Intrinsics.m65267c(this.title, groupCreate.title) && Intrinsics.m65267c(this.avatar_url, groupCreate.avatar_url) && Intrinsics.m65267c(this.conversation_key_version, groupCreate.conversation_key_version);
    }

    public int hashCode() {
        List list = this.member_ids;
        int iHashCode = (list == null ? 0 : list.hashCode()) * 31;
        List list2 = this.admin_ids;
        int iHashCode2 = (iHashCode + (list2 == null ? 0 : list2.hashCode())) * 31;
        String str = this.title;
        int iHashCode3 = (iHashCode2 + (str == null ? 0 : str.hashCode())) * 31;
        String str2 = this.avatar_url;
        int iHashCode4 = (iHashCode3 + (str2 == null ? 0 : str2.hashCode())) * 31;
        String str3 = this.conversation_key_version;
        return iHashCode4 + (str3 != null ? str3.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        List list = this.member_ids;
        List list2 = this.admin_ids;
        String str = this.title;
        String str2 = this.avatar_url;
        String str3 = this.conversation_key_version;
        StringBuilder sb = new StringBuilder("GroupCreate(member_ids=");
        sb.append(list);
        sb.append(", admin_ids=");
        sb.append(list2);
        sb.append(", title=");
        C0026b.m37b(sb, str, ", avatar_url=", str2, ", conversation_key_version=");
        return C0003b.m4b(sb, str3, Separators.RPAREN);
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}